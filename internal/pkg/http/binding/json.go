package binding

import (
	"errors"
	"net/http"

	"github.com/andreyloginov-afk/catalog-service/internal/pkg/http/httph"
	"github.com/go-playground/form/v4"
)

var formDecoder = form.NewDecoder()

type jsonBinding struct{}

type queryBinding struct{}

func (jsonBinding) Name() string {
	return "JSON"
}

func (queryBinding) Name() string {
	return "URL-QUERY"
}

func (jsonBinding) Bind(req *http.Request, obj any) error {
	if req == nil || req.Body == nil {
		return errors.New("invalid request")
	}

	if err := httph.DecodeJSON(req, obj); err != nil {
		return err
	}
	return validate(obj)
}

func (queryBinding) Bind(req *http.Request, obj any) error {
	values := req.URL.Query()

	if err := formDecoder.Decode(obj, values); err != nil {
		return err
	}
	return validate(obj)
}
