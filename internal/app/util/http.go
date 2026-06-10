package util

import (
	"net/http"
	"strings"
)

// возвращает true для маршрутов которые не нужно логировать
func IsFilteredHttpRoute(r *http.Request) bool {
	uri := r.RequestURI

	return strings.Contains(uri, "health") ||
		strings.Contains(uri, "debug") ||
		strings.Contains(uri, "metric")
}
