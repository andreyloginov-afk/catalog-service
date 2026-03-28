package httph

import "net/http"

type Error struct {
	Massage string `json:"error"`
}

func ErrorApply(w http.ResponseWriter, code int, massage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_ = EncodeJSON(w, Error{
		Massage: massage,
	})
}
