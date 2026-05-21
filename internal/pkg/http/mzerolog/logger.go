package mzerolog

import (
	"net/http"
	"strings"
	"time"

	"github.com/andreyloginov-afk/catalog-service/internal/pkg/http/httph"
	"github.com/rs/zerolog"
)

type middleware struct {
	log zerolog.Logger

	fromOptions struct {
		skipper func(r *http.Request) bool
	}
}

// Callback - это основной обработчик middleware
func (m *middleware) Callback(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const (
			tailSucces = " finished with no error"
			tailFail   = " finished (or aborted) with error"
		)
		start := time.Now()
		next.ServeHTTP(w, r)

		err := httph.ErrorGet(r)
		execTime := time.Since(start)

		if m.fromOptions.skipper(r) {
			return
		}

		var mb strings.Builder
		mb.Grow(48 + len(r.RequestURI))
		mb.WriteString(r.Method)
		mb.WriteByte(' ')
		mb.WriteString(r.RequestURI)

		var ev *zerolog.Event
		if err == nil {
			mb.WriteString(tailSucces)
			ev = m.log.Debug()
		} else {
			mb.WriteString(tailFail)
			ev = m.log.Error()
		}

		ev.Err(err)
		ev.Ctx(r.Context())
		ev.Str("exec_time", execTime.String())
		ev.Str("client_ip", r.RemoteAddr)
		ev.Msg(mb.String())
	})
}
