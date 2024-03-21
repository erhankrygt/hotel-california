package apierror

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"hotel-california-backend/internal/localization"
	"net/http"
)

// error codes and names
const (
	CodeInternalServerError = iota + 1
	CodeBadRequestError
	CodeValidationError
	CodeUnauthorizedError
	CodeCouldNotCreateReservationError
	CodeCouldNotChangeReservationCheckInDateError
)

// error names
const (
	NameInternalServerError                       = "InternalServerError"
	NameBadRequestError                           = "BadRequestError"
	NameValidationError                           = "ValidationError"
	NameUnauthorizedError                         = "UnauthorizedError"
	NameCouldNotCreateReservationError            = "CouldNotCreateReservationError"
	NameCouldNotChangeReservationCheckInDateError = "CouldNotChangeReservationCheckInDateError"
)

// compile-time proof of error interface implementation
var _ error = (*APIError)(nil)

// compile-time proofs of localizer interface implementation
var _ localization.Localizer = (*APIError)(nil)

// APIError represents api error
type APIError struct {
	Message             string `json:"message"`
	Name                string `json:"name"`
	Code                int    `json:"code"`
	StatusCode          int    `json:"statusCode"`
	ErrorAction         int    `json:"errorAction,omitempty"`
	BaseError           error  `json:"-"`
	MessageLocalizerKey string `json:"-"`
}

// DefaultInternalServerError represents default internal server error
var DefaultInternalServerError = &APIError{
	Name:                NameInternalServerError,
	Code:                CodeInternalServerError,
	StatusCode:          http.StatusInternalServerError,
	MessageLocalizerKey: "default-internal-server-error-message",
}

// DefaultUnauthorizedError represents default unauthorized error
var DefaultUnauthorizedError = &APIError{
	Name:                NameUnauthorizedError,
	Code:                CodeUnauthorizedError,
	StatusCode:          http.StatusUnauthorized,
	MessageLocalizerKey: "default-unauthorized-error-message",
}

var CouldNotCreateReservation = &APIError{
	Name:                NameCouldNotCreateReservationError,
	Code:                CodeCouldNotCreateReservationError,
	StatusCode:          http.StatusBadRequest,
	MessageLocalizerKey: localization.CouldNotCreateReservation,
}

var CouldNotChangeReservationCheckInDate = &APIError{
	Name:                NameCouldNotChangeReservationCheckInDateError,
	Code:                CodeCouldNotChangeReservationCheckInDateError,
	StatusCode:          http.StatusBadRequest,
	MessageLocalizerKey: localization.CouldNotChangeReservationCheckInDate,
}

// NewBadRequestError returns bad request error
func NewBadRequestError(message error) *APIError {
	return &APIError{
		Message:    message.Error(),
		Name:       NameBadRequestError,
		Code:       CodeBadRequestError,
		StatusCode: http.StatusBadRequest,
	}
}

// NewValidationError returns validation error
func NewValidationError(message string, messageLocalizerKey string) *APIError {
	return &APIError{
		Message:             message,
		Name:                NameValidationError,
		Code:                CodeValidationError,
		StatusCode:          http.StatusBadRequest,
		MessageLocalizerKey: messageLocalizerKey,
	}
}

func (apiErr *APIError) Localize(l *i18n.Localizer) {
	apiErr.Message = localization.Localize(l, apiErr.MessageLocalizerKey, apiErr.Message)
}

func (apiErr *APIError) Error() string {
	return apiErr.Message
}
