package httph

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

func SendError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "aplication/json")
	w.WriteHeader(status)
	// то же самое
	if encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); encodeErr != nil {
		log.Printf("SendError: failed to encode error response: %v", encodeErr)
	}
}
