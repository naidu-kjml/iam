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
	router.GET("/healthcheck", healthcheck)

	// Prevent setting two routes
	if servePath != "/" {
		router.GET(servePath+"healthcheck", healthcheck)
	}

	// App Routes
	router.GET(servePath, sayHello)
	router.GET(servePath+"v1/user", security.AuthWrapper(getOktaUserByEmail(oktaClient), secretManager))
	router.GET(servePath+"user/okta", security.AuthWrapper(getOktaUserByEmail(oktaClient), secretManager))

	router.PanicHandler = panicHandler

	return router
}
