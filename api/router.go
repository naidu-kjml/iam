package api

import (
	"gitlab.skypicker.com/platform/security/iam/security"
	"gitlab.skypicker.com/platform/security/iam/services/okta"
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/julienschmidt/httprouter"
)

// CreateRouter creates a new router instance
func CreateRouter(serviceName string, oktaClient *okta.Client, secretManager security.SecretManager) *httprouter.Router {
	router := httprouter.New(httprouter.WithServiceName(serviceName))

	// Healthcheck routes.
	router.GET("/healthcheck", healthcheck)

	// Hello World Route
	router.GET("/", sayHello)

	// App routes
	router.GET(
		"/v1/user",
		security.AuthWrapper(
			getOktaUserByEmail(oktaClient),
			secretManager,
		),
	)
	router.GET(
		"/v1/teams",
		security.AuthWrapper(
			getTeams(oktaClient),
			secretManager,
		),
	)
	router.GET(
		"/user/okta", security.AuthWrapper(
			addDeprecationWarning(
				getOktaUserByEmail(oktaClient),
			),
			secretManager,
		),
	)

	router.PanicHandler = panicHandler

	return router
}
