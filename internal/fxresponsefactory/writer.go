package fxresponsefactory

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ecumenos/ecumenos/internal/fxtypes"
	"github.com/ecumenos/ecumenos/internal/toolkit/contextutils"
	"github.com/ecumenos/ecumenos/internal/toolkit/httputils"
	"github.com/ecumenos/ecumenos/internal/toolkit/timeutils"
	"go.uber.org/zap"
)

//go:generate mockery --name=Writer

type Writer interface {
	SetLogger(logger *zap.Logger)
	WriteSuccess(ctx context.Context, payload interface{}, opts ...ResponseBuildOption) error
	WriteFail(ctx context.Context, data interface{}, opts ...ResponseBuildOption) error
	WriteError(ctx context.Context, msg string, cause error, opts ...ResponseBuildOption) error
}

type writer struct {
	rw         http.ResponseWriter
	l          *zap.Logger
	writeLogs  bool
	appVersion fxtypes.Version
}

func NewWriter(l *zap.Logger, rw http.ResponseWriter, appVersion fxtypes.Version, writeLogs bool) Writer {
	return &writer{
		rw:         rw,
		l:          l,
		writeLogs:  writeLogs,
		appVersion: appVersion,
	}
}

func (w *writer) write(payload interface{}) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if _, err := w.rw.Write(b); err != nil {
		return err
	}

	return nil
}

func (w *writer) writeHeaders(headers map[string]string, statusCode int) {
	w.rw.Header().Set("Content-Type", "application/json")
	for key, value := range headers {
		w.rw.Header().Set(key, value)
	}
	w.rw.WriteHeader(statusCode)
}

func (w *writer) SetLogger(logger *zap.Logger) {
	w.l = logger
}

type Status string

const (
	SuccessStatus Status = "success"
	FailureStatus Status = "failure"
	ErrorStatus   Status = "error"
)

type SuccessResp[T interface{}] struct {
	Status Status `json:"status"`
	Data   T      `json:"data"`
}

func (w *writer) WriteSuccess(ctx context.Context, payload interface{}, opts ...ResponseBuildOption) error {
	headers, err := w.getHeaders(ctx)
	if err != nil {
		return err
	}
	rb := &responseBuilder{
		httpStatusCode: http.StatusOK,
		data:           payload,
		l:              w.l,
	}
	for _, opt := range opts {
		opt(rb)
	}
	if rb.httpStatusCode < http.StatusOK || rb.httpStatusCode > 299 {
		return fmt.Errorf("success response must have status code in range 200..299 (status code = %v)", rb.httpStatusCode)
	}
	if w.writeLogs {
		w.l.Info("responding success response", zap.Any("data", payload), zap.Int("status_code", rb.httpStatusCode))
	}

	w.writeHeaders(headers, rb.httpStatusCode)
	return w.write(&SuccessResp[interface{}]{
		Data:   rb.data,
		Status: SuccessStatus,
	})
}

type FailureResp[T interface{}] struct {
	Status  Status `json:"status"`
	Data    T      `json:"data"`
	Message string `json:"message"`
}

func (w *writer) WriteFail(ctx context.Context, data interface{}, opts ...ResponseBuildOption) error {
	headers, err := w.getHeaders(ctx)
	if err != nil {
		return err
	}
	rb := &responseBuilder{
		httpStatusCode: http.StatusBadRequest,
		data:           data,
		l:              w.l,
	}
	for _, opt := range opts {
		opt(rb)
	}
	if rb.httpStatusCode < http.StatusBadRequest || rb.httpStatusCode > 499 {
		return fmt.Errorf("fail response must have status code in range 400..499 (status code = %v)", rb.httpStatusCode)
	}
	if w.writeLogs {
		w.l.Info("responding fail response", zap.Any("data", rb.data), zap.Error(rb.cause),
			zap.Int("status_code", rb.httpStatusCode), zap.String("msg", rb.message))
	}

	w.writeHeaders(headers, http.StatusBadRequest)
	return w.write(&FailureResp[interface{}]{
		Data:    rb.data,
		Message: rb.message,
		Status:  FailureStatus,
	})
}

type ErrorResp struct {
	Status  Status `json:"status"`
	Message string `json:"message"`
}

func (w *writer) WriteError(ctx context.Context, msg string, cause error, opts ...ResponseBuildOption) error {
	headers, err := w.getHeaders(ctx)
	if err != nil {
		return err
	}
	rb := &responseBuilder{
		httpStatusCode: http.StatusInternalServerError,
		message:        msg,
		cause:          cause,
		l:              w.l,
	}
	for _, opt := range opts {
		opt(rb)
	}
	if rb.httpStatusCode < http.StatusInternalServerError || rb.httpStatusCode > 599 {
		return fmt.Errorf("error response must have status code in range 500..599 (status code = %v)", rb.httpStatusCode)
	}
	if w.writeLogs {
		w.l.Info("responding error response", zap.Error(rb.cause), zap.String("msg", rb.message), zap.Int("status_code", rb.httpStatusCode))
	}

	w.writeHeaders(headers, rb.httpStatusCode)
	return w.write(&ErrorResp{
		Message: rb.message,
		Status:  ErrorStatus,
	})
}

type responseBuilder struct {
	httpStatusCode int
	message        string
	cause          error
	data           interface{}
	l              *zap.Logger
}

type ResponseBuildOption func(b *responseBuilder)

func WithHTTPStatusCode(code int) ResponseBuildOption {
	return func(b *responseBuilder) {
		b.httpStatusCode = code
	}
}

func WithMessage(msg string) ResponseBuildOption {
	return func(b *responseBuilder) {
		b.message = msg
	}
}

func WithCause(err error) ResponseBuildOption {
	return func(b *responseBuilder) {
		b.cause = err
	}
}

func WithData(data interface{}) ResponseBuildOption {
	return func(b *responseBuilder) {
		b.data = data
	}
}

func WithLogger(l *zap.Logger) ResponseBuildOption {
	return func(b *responseBuilder) {
		b.l = l
	}
}

func (w *writer) getHeaders(ctx context.Context) (map[string]string, error) {
	duration, err := httputils.GetRequestDuration(ctx)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"X-Request-Id":              contextutils.GetRequestID(ctx),
		"X-Timestamp":               timeutils.TimeToString(time.Now()),
		"X-Request-Duration":        fmt.Sprint(duration.Milliseconds()),
		"X-App-Version":             string(w.appVersion),
		"Cache-Control":             "no-cache, no-store, max-age=0, must-revalidate",
		"Pragma":                    "no-cache",
		"Expires":                   "0",
		"X-Content-Type-Options":    "nosniff",
		"Strict-Transport-Security": "max-age=31536000 ; includeSubDomains",
		"X-Frame-Options":           "DENY",
		"X-XSS-Protection":          "0",
		"Content-Type":              "application/json",
	}, nil
}
