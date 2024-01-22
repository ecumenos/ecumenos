package httputils

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"
)

func DecodeBody[Out any](l *zap.Logger, r *http.Request) (*Out, error) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			l.Error("failed to close body", zap.Error(err))
		}
	}(r.Body)

	var body Out
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}
