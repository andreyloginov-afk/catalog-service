package httph

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

func SendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "aplication/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// спецом залогировал ошибку линтер ругается
		log.Printf("SendJSON: failed to encode response: %v", err)
	}
}

func SendEmpty(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func SendError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "aplication/json")
	w.WriteHeader(status)
	// то же самое
	if encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); encodeErr != nil {
		log.Printf("SendError: failed to encode error response: %v", encodeErr)
	}
}
