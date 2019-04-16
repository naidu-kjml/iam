package api

import (
	"gitlab.skypicker.com/platform/security/iam/security"
	"gitlab.skypicker.com/platform/security/iam/services/okta"
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/julienschmidt/httprouter"
)

// CreateRouter creates a new router instance
func CreateRouter(serviceName, servePath string, oktaClient *okta.Client, secretManager security.SecretManager) *httprouter.Router {
	router := httprouter.New(httprouter.WithServiceName(serviceName))

	// Healthcheck routes. Exposed on both /healthcheck and /servePath/healthcheck to allow easier k8s set up
	router.GET("/healthcheck", Healthcheck)

	// Prevent setting two routes
	if servePath != "/" {
		router.GET(servePath+"healthcheck", Healthcheck)
	}

	// App Routes
	router.GET(servePath, SayHello)
	router.GET(servePath+"user/okta", security.AuthWrapper(GetOktaUserByEmail(oktaClient), secretManager))

	router.PanicHandler = PanicHandler

	return router
}
