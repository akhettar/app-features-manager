package model

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"time"
)

const (
	AppVersion  = "version"
	AppPlatform = "platform"
)

// App Status enum
type Status string

// App Platform
type Platform string

const (
	// App Status
	Latest      Status = "latest"
	Deprecated         = "deprecated"
	Unsupported        = "unsupported"
	Supported          = "supported"

	// App platform
	Ios        Platform = "ios"
	Android             = "android"
	Windows             = "windows"
	Blackberry          = "blackberry"
)

// Value returns the value of the app status enum or an error if there is no match
func (s Status) Value() (string, error) {
	switch s {
	case Deprecated, Unsupported, Supported, Latest:
		return string(s), nil
	}
	return "", errors.New(fmt.Sprintf("Supported values:%s,%s,%s,%s", Deprecated, Unsupported, Supported, Latest))
}

// Value returns the value of the app platform enum or an error if there is no match
func (p Platform) Value() (string, error) {
	switch p {
	case Ios, Android, Windows, Blackberry:
		return string(p), nil
	}
	return "", errors.New(fmt.Sprintf("Supported values:%s,%s,%s,%s", Ios, Android, Windows, Blackberry))
}

// ReleaseDAO instance of the app status to be stored in the data store
type ReleaseDAO struct {
	Version  string    `bson:"version"`
	Status   string    `bson:"status"`
	Platform string    `bson:"platform"`
	Released time.Time `bson:"released"`
}

// ReleaseRequest is the payload to releasing the app version
type (
	ReleaseRequest struct {
		Version  string `json:"version" validate:"required"`
		Platform string `json:"platform" validate:"required"`
		Status   string `json:"status"`
	}

	ReleaseRequestValidator struct {
		Validator *validator.Validate
	}
)

func (cv *ReleaseRequestValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

// ReleaseRequestStructLevelValidation contains custom struct level validations that don't always
// make sense at the field validation level.
func ReleaseRequestStructLevelValidation(sl validator.StructLevel) {

	req := sl.Current().Interface().(ReleaseRequest)

	// validate the status
	if _, e := Status(req.Status).Value(); e != nil {
		sl.ReportError(req.Status, "status", "Status", "", "")
	}
	if _, e := Platform(req.Platform).Value(); e != nil {
		sl.ReportError(req.Platform, "platform", "Platform", "", "")
	}
}

// ReleaseResponse is the query app status response
type ReleaseResponse struct {
	Status string          `json:"status"`
	Flags  map[string]bool `json:"flags"`
}

// ErrorResponse a generic error response
type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:code`
}

// EmptyBody for version not found
type EmptyBody struct{}
