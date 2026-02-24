package health

import (
	"net/http"

	rhandler "github.com/andreyloginov-afk/catalog-service/internal/app/handler"
)

type handler struct{}

func NewHandler() rhandler.Health {
	return &handler{}
}

func (h *handler) LastCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
