package health

import (
	"net/http"

	rhandler "github.com/andreyloginov-afk/catalog-service/internal/app/handler"
	"github.com/rs/zerolog/log"
)

type handler struct{}

func NewHandler() rhandler.Health {
	return &handler{}
}

func (h *handler) LastCheck(w http.ResponseWriter, r *http.Request) {
	// залогировал ошибку
	if _, err := w.Write([]byte("ok")); err != nil {
		log.Error().Err(err).Msg("failed to write health check response")
	}
}
