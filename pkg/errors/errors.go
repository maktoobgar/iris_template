package errors

import (
	"fmt"
	"net/http"
	"service/pkg/colors"
	"strings"
)

type (
	serverError struct {
		code    int64
		action  int64
		message string
		errMsg  string
		errors  any
	}
)

const (
	_ int64 = iota
	// BadRequest 400
	InvalidStatus
	// NotFound 404
	NotFoundStatus
	// Unauthorized 401
	UnauthorizedStatus
	// InternalServerError 500
	UnexpectedStatus
	// MethodNotAllowed 405
	MethodNotAllowedStatus
	// Forbidden 403
	ForbiddenStatus
	// Timeout 408
	TimeoutStatus
	// Service Unavailable 503
	ServiceUnavailable
)

const (
	// Do nothing
	DoNothing int64 = iota
	// SignIn in again
	ReSignIn
	// Report the problem
	Report
	// Correct sent data and request again
	Resend
	// Try later
	TryLater
	// Refresh your access token
	RefreshToken
)

var (
	httpErrors = map[int64]int{
		InvalidStatus:          http.StatusBadRequest,
		NotFoundStatus:         http.StatusNotFound,
		UnauthorizedStatus:     http.StatusUnauthorized,
		UnexpectedStatus:       http.StatusInternalServerError,
		MethodNotAllowedStatus: http.StatusMethodNotAllowed,
		ForbiddenStatus:        http.StatusForbidden,
		TimeoutStatus:          http.StatusRequestTimeout,
		ServiceUnavailable:     http.StatusServiceUnavailable,
	}
)

func (e serverError) Error() string {
	errCode := httpErrors[e.code]
	if errCode >= 500 {
		return fmt.Sprintf("%sCode: %d - Action: %d - Message: %s - Error: %s%s", colors.Red, errCode, e.action, e.message, e.errMsg, colors.Reset)
	}
	message := fmt.Sprintf("%sCode: %d - Action: %d - Message: %s - Error: %s%s", colors.Orange, errCode, e.action, e.message, e.errMsg, colors.Reset)
	if e.errMsg == "" {
		message = fmt.Sprintf("%sCode: %d - Action: %d - Message: %s%s", colors.Green, errCode, e.action, e.message, colors.Reset)
	}
	return message
}

// Returns httpErrorCode, message and action of it
func HttpError(err error) (code int, action int, message string, errMsg string, errors any) {
	code = http.StatusInternalServerError
	action = int(Report)
	errMsg = err.Error()
	message = errMsg
	errors = nil

	if er, ok := err.(serverError); ok {
		code = httpErrors[er.code]
		action = int(er.action)
		errMsg = er.errMsg
		message = er.message
		errors = er.errors
	}

	return
}

func IsServerError(err error) bool {
	if _, ok := err.(serverError); ok {
		return true
	}
	return false
}

// Creates a new error
func New(code int64, action int64, message string, errMsg string, errors ...any) error {
	var errs any = nil
	if len(errors) != 0 {
		errs = errors[0]
	}
	return serverError{
		code:    code,
		action:  action,
		errMsg:  errMsg,
		message: message,
		errors:  errs,
	}
}

func IsContextDeadlineExceeded(err error) bool {
	return strings.Contains(err.Error(), "context deadline exceeded") || strings.Contains(err.Error(), "context canceled")
}
