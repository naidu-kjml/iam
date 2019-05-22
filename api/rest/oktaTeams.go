package rest

import (
	"log"
	"net/http"

	"github.com/getsentry/raven-go"
	"github.com/julienschmidt/httprouter"
)

type teamsGetter interface {
	GetTeams() (map[string]int, error)
}

func getTeams(client teamsGetter) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		teams, err := client.GetTeams()
		if err != nil {
			http.Error(w, "Service unavailable", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		je := json.NewEncoder(w)
		if err := je.Encode(teams); err != nil {
			log.Println("[ERROR]", err.Error())
			raven.CaptureError(err, nil)
		}
	}
}
