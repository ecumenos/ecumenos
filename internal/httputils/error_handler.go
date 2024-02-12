package httputils

import (
	"net/http"

	"github.com/ecumenos/ecumenos/internal/fxresponsefactory"
)

func DefaultErrorHandlerFactory(rf fxresponsefactory.Factory) func(rw http.ResponseWriter, r *http.Request, err error) {
	return func(rw http.ResponseWriter, r *http.Request, err error) {
		opts := []fxresponsefactory.ResponseBuildOption{
			fxresponsefactory.WithHTTPStatusCode(http.StatusBadRequest),
			fxresponsefactory.WithCause(err),
		}
		_ = rf.NewWriter(rw).WriteFail(r.Context(), nil, opts...) //nolint:errcheck
	}
}
