package rest

import (
	"net/http"
	"strings"

	"gitlab.skypicker.com/platform/security/iam/security"
	"gitlab.skypicker.com/platform/security/iam/security/permissions"
	"gitlab.skypicker.com/platform/security/iam/security/secrets"
	"gitlab.skypicker.com/platform/security/iam/services/okta"
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/julienschmidt/httprouter"
)

const wellKnownFolder string = ".well-known"

type metricService interface {
	// Incr increments by 1 a metric identified by name.
	// tags should be in format name:value and can be created with Tag function to escape the values
	Incr(string, ...string)
}

// CreateRouter creates a new router instance
func CreateRouter(
	serviceName string,
	oktaClient *okta.Client,
	permissionManager permissions.PermissionManager,
	secretManager secrets.SecretManager,
	metricClient metricService) *httprouter.Router {
	router := httprouter.New(httprouter.WithServiceName(serviceName))

	router.Handler(
		"GET",
		"/"+wellKnownFolder+"/*filepath",
		DisableDirectoryListingHandler(
			http.StripPrefix("/"+wellKnownFolder+"/", http.FileServer(http.Dir(wellKnownFolder))),
		),
	)

	// Healthcheck routes.
	router.GET("/healthcheck", healthcheck)

	// Hello World Route
	router.GET("/", sayHello)

	// App routes
	router.GET(
		"/v1/user",
		security.AuthWrapper(
			getOktaUserByEmail(oktaClient, permissionManager),
			secretManager,
			metricClient,
		),
	)
	router.GET(
		"/v1/teams",
		security.AuthWrapper(
			getTeams(oktaClient),
			secretManager,
			metricClient,
		),
	)
	router.GET(
		"/v1/groups",
		security.AuthWrapper(
			getGroups(oktaClient),
			secretManager,
			metricClient,
		),
	)
	router.GET(
		"/user/okta", security.AuthWrapper(
			addDeprecationWarning(
				getOktaUserByEmail(oktaClient, permissionManager),
			),
			secretManager,
			metricClient,
		),
	)

	router.PanicHandler = panicHandler

	return router
}

// DisableDirectoryListingHandler prevents directory listings to be returned to the user-agent.
func DisableDirectoryListingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}
