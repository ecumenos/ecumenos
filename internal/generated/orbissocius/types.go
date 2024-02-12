// Package orbissocius provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package orbissocius

import (
	"encoding/json"

	"github.com/oapi-codegen/runtime"
)

const (
	BearerAuthScopes = "bearerAuth.Scopes"
)

// Defines values for ErrorResponseStatus.
const (
	ErrorResponseStatusError ErrorResponseStatus = "error"
)

// Defines values for FailResponseStatus.
const (
	Fail FailResponseStatus = "fail"
)

// Defines values for SuccessResponseStatus.
const (
	SuccessResponseStatusSuccess SuccessResponseStatus = "success"
)

// ErrorResponseBody defines model for ErrorResponseBody.
type ErrorResponseBody struct {
	// Message A meaningful, end-user-readable message, explaining what went wrong.
	Message string              `json:"message"`
	Status  ErrorResponseStatus `json:"status"`
}

// ErrorResponseStatus defines model for ErrorResponseStatus.
type ErrorResponseStatus string

// FailResponseStatus defines model for FailResponseStatus.
type FailResponseStatus string

// FailureResponseBody defines model for FailureResponseBody.
type FailureResponseBody struct {
	Data *map[string]interface{} `json:"data"`

	// Message A meaningful, end-user-readable message, explaining what went wrong.
	Message *string            `json:"message,omitempty"`
	Status  FailResponseStatus `json:"status"`
}

// GetHealthData defines model for GetHealthData.
type GetHealthData struct {
	// Ok The OK is true if all is okay.
	Ok bool `json:"ok"`
}

// GetInfoData defines model for GetInfoData.
type GetInfoData struct {
	Deps *map[string]interface{} `json:"deps"`

	// Name The Name is the service name.
	Name string `json:"name"`

	// Version The Version is the service semver version.
	Version string `json:"version"`
}

// JSendResponseArray defines model for JSendResponseArray.
type JSendResponseArray struct {
	Data   *[]map[string]interface{} `json:"data"`
	Status SuccessResponseStatus     `json:"status"`
}

// JSendResponseObject defines model for JSendResponseObject.
type JSendResponseObject struct {
	Data   *map[string]interface{} `json:"data"`
	Status SuccessResponseStatus   `json:"status"`
}

// ResponseStatus defines model for ResponseStatus.
type ResponseStatus struct {
	union json.RawMessage
}

// SuccessResponseStatus defines model for SuccessResponseStatus.
type SuccessResponseStatus string

// BadRequest defines model for BadRequest.
type BadRequest = FailureResponseBody

// Error defines model for Error.
type Error = ErrorResponseBody

// Failure defines model for Failure.
type Failure = FailureResponseBody

// Forbidden defines model for Forbidden.
type Forbidden = FailureResponseBody

// NotAuthorized defines model for NotAuthorized.
type NotAuthorized = FailureResponseBody

// NotFound defines model for NotFound.
type NotFound = FailureResponseBody

// Success defines model for Success.
type Success = JSendResponseObject

// AsSuccessResponseStatus returns the union data inside the ResponseStatus as a SuccessResponseStatus
func (t ResponseStatus) AsSuccessResponseStatus() (SuccessResponseStatus, error) {
	var body SuccessResponseStatus
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromSuccessResponseStatus overwrites any union data inside the ResponseStatus as the provided SuccessResponseStatus
func (t *ResponseStatus) FromSuccessResponseStatus(v SuccessResponseStatus) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeSuccessResponseStatus performs a merge with any union data inside the ResponseStatus, using the provided SuccessResponseStatus
func (t *ResponseStatus) MergeSuccessResponseStatus(v SuccessResponseStatus) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

// AsFailResponseStatus returns the union data inside the ResponseStatus as a FailResponseStatus
func (t ResponseStatus) AsFailResponseStatus() (FailResponseStatus, error) {
	var body FailResponseStatus
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromFailResponseStatus overwrites any union data inside the ResponseStatus as the provided FailResponseStatus
func (t *ResponseStatus) FromFailResponseStatus(v FailResponseStatus) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeFailResponseStatus performs a merge with any union data inside the ResponseStatus, using the provided FailResponseStatus
func (t *ResponseStatus) MergeFailResponseStatus(v FailResponseStatus) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

// AsErrorResponseStatus returns the union data inside the ResponseStatus as a ErrorResponseStatus
func (t ResponseStatus) AsErrorResponseStatus() (ErrorResponseStatus, error) {
	var body ErrorResponseStatus
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromErrorResponseStatus overwrites any union data inside the ResponseStatus as the provided ErrorResponseStatus
func (t *ResponseStatus) FromErrorResponseStatus(v ErrorResponseStatus) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeErrorResponseStatus performs a merge with any union data inside the ResponseStatus, using the provided ErrorResponseStatus
func (t *ResponseStatus) MergeErrorResponseStatus(v ErrorResponseStatus) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

func (t ResponseStatus) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *ResponseStatus) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}
