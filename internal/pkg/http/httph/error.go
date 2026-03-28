package httph

import "net/http"

type Error struct {
	Message string `json:"error"`
}

func ErrorApply(w http.ResponseWriter, code int, massage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_ = EncodeJSON(w, Error{
		Message: massage,
	})
}
