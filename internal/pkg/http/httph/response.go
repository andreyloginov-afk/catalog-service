package httph

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"
)

type httpCoder interface {
	HTTPStatus() int
	Error() string
}

func SendError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); encodeErr != nil {
		log.Error().
			Err(encodeErr).
			Int("status", status).
			Msg("failed to encode error response")
	}
}

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	ErrorApply(r, err)

	var hc httpCoder
	if errors.As(err, &hc) {
		status := hc.HTTPStatus()
		ErrorApplyStatusCode(r, status)
		SendError(w, status, hc)
		return
	}

	ErrorApplyStatusCode(r, http.StatusInternalServerError)
	SendError(w, http.StatusInternalServerError, err)
}
