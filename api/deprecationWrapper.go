package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func addDeprecationWarning(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Deprecated", "true")

		h(w, r, ps)
	}
}
