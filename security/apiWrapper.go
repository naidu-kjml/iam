package security

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/julienschmidt/httprouter"
	"gitlab.skypicker.com/go/packages/useragent"
	"gitlab.skypicker.com/platform/security/iam/shared"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// AuthWrapper wraps a router to validate the authentication token
func AuthWrapper(h httprouter.Handle, secretManager SecretManager) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		err := checkAuth(r, secretManager)
		if err != nil {
			if apiErr, ok := err.(shared.APIError); ok {
				http.Error(w, apiErr.Message, apiErr.Code)
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}

			log.Println("[ERROR]", err.Error())
			return
		}

		// Delegate request to the given handle
		h(w, r, ps)
	}
}

const (
	serviceNamePattern string = `^[\w\s-]+`
)

var serviceNameRe = regexp.MustCompile(serviceNamePattern)

// It's important to add $ in order to match the whole string
var checkServiceNameRe = regexp.MustCompile(serviceNamePattern + "$")

func checkServiceName(service string) error {
	safe := checkServiceNameRe.MatchString(service)
	if !safe {
		return errors.New("service name has to match " + serviceNamePattern + "$")
	}

	return nil
}

// GetServiceName returns the service name based on the given user agent.
func GetServiceName(incomingUserAgent string) (string, error) {
	ua, err := useragent.Parse(incomingUserAgent)

	if err == nil {
		return strings.ToUpper(ua.Name), nil
	}

	// Log is temp. This should be pushed to Datadog when possible.
	log.Printf("User agent [%v] failed: [%v]", incomingUserAgent, err)

	// This block should be removed after all services adhere to RFC 22.
	service := serviceNameRe.FindString(incomingUserAgent)
	if service == "" {
		return "", errors.New("no service found")
	}

	service = strings.ToUpper(service)
	service = strings.ReplaceAll(service, " ", "_")

	return service, nil
}

// checkAuth checks if user has proper token + user agent
func checkAuth(r *http.Request, secretManager SecretManager) error {
	var requestToken = r.Header.Get("Authorization")
	var userAgent = r.Header.Get("User-Agent")

	service, err := GetServiceName(userAgent)
	if err != nil {
		return shared.APIError{Message: "User-Agent header mandatory", Code: 401}
	}

	if requestToken == "" {
		return shared.APIError{Message: "Authorization header with token is mandatory", Code: 401}
	}

	if span, ok := tracer.SpanFromContext(r.Context()); ok {
		span.SetTag("user-agent", userAgent)
		span.SetTag("service-name", service)
	}

	token, err := secretManager.GetAppToken(service)

	if err != nil {
		return shared.APIError{Message: "Missing token", Code: 401}
	}

	if token != requestToken {
		return shared.APIError{Message: "Incorrect token", Code: 401}
	}

	return nil
}
