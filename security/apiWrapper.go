package security

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/julienschmidt/httprouter"
	"gitlab.skypicker.com/go/packages/useragent"
	"gitlab.skypicker.com/platform/security/iam/api"
	"gitlab.skypicker.com/platform/security/iam/monitoring"
	"gitlab.skypicker.com/platform/security/iam/security/secrets"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type metricService interface {
	// Incr increments by 1 a metric identified by name.
	// tags should be in format name:value and can be created with Tag function to escape the values
	Incr(string, ...string)
}

// AuthWrapper wraps a router to validate the authentication token
func AuthWrapper(h httprouter.Handle, secretManager secrets.SecretManager, metricClient metricService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		err := checkAuth(r, secretManager, metricClient)
		if err != nil {
			if apiErr, ok := err.(api.Error); ok {
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

// CheckServiceName returns if the given service name contains expected characters only
func CheckServiceName(service string) error {
	safe := checkServiceNameRe.MatchString(service)
	if !safe {
		return errors.New("service name has to match " + serviceNamePattern + "$")
	}

	return nil
}

// Service holds the requesting service's name and environment
type Service struct {
	Name        string
	Environment string
}

// GetService returns the service name and environment based on the given user agent
func GetService(incomingUserAgent string) (Service, error) {
	ua, err := useragent.Parse(incomingUserAgent)

	if err == nil {
		return Service{ua.Name, ua.Environment}, nil
	}
	// Log is temp. This should be pushed to Datadog when possible
	log.Printf("User agent [%v] failed: [%v]", incomingUserAgent, err)

	// This block should be removed after all services adhere to RFC 22
	service := serviceNameRe.FindString(incomingUserAgent)
	if service == "" {
		return Service{}, errors.New("no service found")
	}

	service = strings.ToUpper(service)
	service = strings.ReplaceAll(service, " ", "_")

	return Service{service, ""}, nil
}

// checkAuth checks if user has proper token + user agent
func checkAuth(r *http.Request, secretManager secrets.SecretManager, metricClient metricService) error {
	requestToken := getToken(r.Header.Get("Authorization"))
	userAgent := r.Header.Get("User-Agent")

	service, err := GetService(userAgent)
	if err != nil {
		return api.Error{Message: "User-Agent header mandatory", Code: 401}
	}

	if requestToken == "" {
		return api.Error{Message: "Authorization header with token is mandatory", Code: 401}
	}

	if span, ok := tracer.SpanFromContext(r.Context()); ok {
		span.SetTag("user-agent", userAgent)
		span.SetTag("service-name", service.Name)
	}

	token, err := secretManager.GetAppToken(service.Name, service.Environment)
	if err != nil {
		return api.Error{Message: "Unauthorized: " + err.Error(), Code: 401}
	}

	if token != requestToken {
		return api.Error{Message: "Unauthorized: incorrect token", Code: 401}
	}
	// Track old authentication format
	if !strings.Contains(r.Header.Get("Authorization"), "Bearer") {
		metricClient.Incr(
			"incoming.old-authentication",
			monitoring.Tag("service-name", service.Name),
			monitoring.Tag("service-environment", service.Environment),
		)
	}
	metricClient.Incr(
		"incoming.requests",
		monitoring.Tag("service-name", service.Name),
		monitoring.Tag("service-environment", service.Environment),
	)

	return nil
}

func getToken(authorization string) string {
	return strings.Replace(authorization, "Bearer ", "", 1)
}
