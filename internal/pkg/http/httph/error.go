package httph

import (
	"context"
	"net/http"
)

type (
	contextKeyError   struct{}
	contextValueError struct {
		err        error
		statusCode int
	}
)

func errorPrepare(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKeyError{}, &contextValueError{})
}

func errorApply(ctx context.Context, err error) {
	if val := ctx.Value(contextKeyError{}); val != nil {
		if v, ok := val.(*contextValueError); ok {
			v.err = err
		}
	}
}

func errorGet(ctx context.Context) error {
	// получаем *contextValueError из контекста
	if val := ctx.Value(contextKeyError{}); val != nil {
		if v, ok := val.(*contextValueError); ok {
			return v.err
		}
	}
	return nil
}

func errorApplyStatusCode(ctx context.Context, statusCode int) {
	// тут аналогично errorApply но для statusCode
	if val := ctx.Value(contextKeyError{}); val != nil {
		if v, ok := val.(*contextValueError); ok {
			v.statusCode = statusCode
		}
	}
}

func errorGetStatusCode(ctx context.Context) int {
	// Аналогично errorGet, но возвращаем statuscode или 0
	if val := ctx.Value(contextKeyError{}); val != nil {
		if v, ok := val.(*contextValueError); ok {
			return v.statusCode
		}
	}
	return 0
}

// публичные обертки для работы с http.Request
func ErrorPrepare(r *http.Request) *http.Request {
	return r.WithContext(errorPrepare(r.Context()))
}

func ErrorApply(r *http.Request, err error) {
	errorApply(r.Context(), err)
}

func ErrorGet(r *http.Request) error {
	return errorGet(r.Context())
}

func ErrorApplyStatusCode(r *http.Request, statusCode int) {
	errorApplyStatusCode(r.Context(), statusCode)
}

func ErrorGetStatusCode(r *http.Request) int {
	return errorGetStatusCode(r.Context())
}

type Middlewear = func(http.Handler) http.Handler

func NewErrorMiddlewear() Middlewear {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, ErrorPrepare(r))
		})
	}
}
