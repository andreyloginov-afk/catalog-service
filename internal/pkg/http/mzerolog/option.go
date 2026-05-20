package mzerolog

import (
	"net/http"

	"github.com/andreyloginov-afk/catalog-service/internal/pkg/http/httph"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Option = func(m *middleware)

func WithLogger(l zerolog.Logger) Option {
	return func(m *middleware) {
		m.log = l

	}
}

func WithSkipper(skipper func(r *http.Request) bool) Option {
	return func(m *middleware) {
		if skipper != nil {
			m.fromOptions.skipper = skipper
		}
	}
}

func NewMiddleware(opts ...Option) httph.Middlewear {
	m := middleware{
		log: log.Logger,
	}
	m.fromOptions.skipper = defaultSkipper

	for _, opt := range opts {
		opt(&m)
	}

	return m.Callback

}

func defaultSkipper(_ *http.Request) bool {
	return false
}
