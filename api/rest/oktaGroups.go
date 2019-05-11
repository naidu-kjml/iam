package api

import (
	"log"
	"net/http"

	"github.com/getsentry/raven-go"
	"github.com/go-redis/redis"
	"github.com/julienschmidt/httprouter"
	"gitlab.skypicker.com/platform/security/iam/services/okta"
)

type groupsGetter interface {
	GetGroups() ([]okta.Group, error)
}

func getGroups(client groupsGetter) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		groups, err := client.GetGroups()
		if err == redis.Nil {
			// No value available for groups yet
			w.Header().Add("Retry-After", "30")
			http.Error(w, "Groups not loaded yet, try later", http.StatusServiceUnavailable)
			return
		}
		if err != nil {
			http.Error(w, "Service unavailable", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		je := json.NewEncoder(w)
		if err := je.Encode(groups); err != nil {
			log.Println("[ERROR]", err.Error())
			raven.CaptureError(err, nil)
		}
	}
}
