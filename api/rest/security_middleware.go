package rest

import (
	"log"
	"net/http"

	"strings"

	"github.com/julienschmidt/httprouter"

	"gitlab.skypicker.com/platform/security/iam/api"
	"gitlab.skypicker.com/platform/security/iam/monitoring"
	"gitlab.skypicker.com/platform/security/iam/security"
	"gitlab.skypicker.com/platform/security/iam/security/secrets"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

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

// checkAuth checks if user has proper token + user agent
func checkAuth(r *http.Request, secretManager secrets.SecretManager, metricClient metricService) error {
	requestToken := security.GetToken(r.Header.Get("Authorization"))
	userAgent := r.Header.Get("User-Agent")

	service, err := security.GetService(userAgent)
	if err != nil {
		return api.Error{Message: err.Error(), Code: 401}
	}

	if requestToken == "" {
		return api.Error{Message: "Authorization header with token is mandatory", Code: 401}
	}

	if span, ok := tracer.SpanFromContext(r.Context()); ok {
		span.SetTag("user-agent", userAgent)
		span.SetTag("service-name", service.Name)
	}

	tokenErr := security.VerifyToken(secretManager, service, requestToken)

	if tokenErr != nil {
		return api.Error{Message: "Unauthorized: " + tokenErr.Error(), Code: 401}
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
