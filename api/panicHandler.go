package api

import (
	"log"
	"net/http"

	"github.com/getsentry/raven-go"
	"gitlab.skypicker.com/platform/security/iam/shared"
)

// PanicHandler prevents the server crashing due to unhandled panics
func PanicHandler(w http.ResponseWriter, r *http.Request, err interface{}) {
	apiError, castingOk := err.(shared.APIError)
	if castingOk {
		raven.CaptureError(apiError, nil)
		http.Error(w, apiError.Message, apiError.Code)
		log.Println("[ERROR]", apiError)
		return
	}

	if errorType, castingOk := err.(error); castingOk {
		raven.CaptureError(errorType, nil)
	}
	log.Panic(err)
}
