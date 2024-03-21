package hotelcalifornia

import (
	"context"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	apierror "hotel-california-backend/internal/api-error"
)

// Service defines behaviors of account service
type Service interface {
	Health(context.Context, HealthRequest) HealthResponse
	SignIn(context.Context, SignInRequest) SignInResponse
	CreateReservation(context.Context, CreateReservationRequest) CreateReservationResponse
	UpdateReservation(context.Context, UpdateReservationRequest) UpdateReservationResponse
	FindReservation(context.Context, FindReservationRequest) FindReservationResponse
	FindReservations(context.Context, FindReservationsRequest) FindReservationsResponse
}

// Request defines behaviors of request
type Request interface {
	SetIPAddress(ipAddress string)
}

// Response defines behaviors of response
type Response interface {
	Localize(*i18n.Localizer) interface{}
	APIError() error
}

// compile-time proofs of request interface implementation
var (
	_ Request = (*HealthRequest)(nil)
	_ Request = (*SignInRequest)(nil)
	_ Request = (*CreateReservationRequest)(nil)
	_ Request = (*UpdateReservationRequest)(nil)
	_ Request = (*FindReservationRequest)(nil)
	_ Request = (*FindReservationsRequest)(nil)
)

// compile-time proofs of response interface implementation
var (
	_ Response = (*HealthResponse)(nil)
	_ Response = (*SignInResponse)(nil)
	_ Response = (*CreateReservationResponse)(nil)
	_ Response = (*UpdateReservationResponse)(nil)
	_ Response = (*FindReservationResponse)(nil)
	_ Response = (*FindReservationsResponse)(nil)
)

type Header struct {
	AcceptLanguage string `header:"Accept-Language" json:"Accept-Language"`
	Token          string `header:"token" json:"token"`
}

type HealthRequest struct {
	IPAddress string `json:"-"`
}

type HealthResponse struct {
	Ping string
}

// sign in models
type (
	SignInRequest struct {
		IPAddress string `json:"-"`
		UserName  string `json:"userName" validate:"required"`
		Password  string `json:"password" validate:"required"`
	}
	// swagger:response SignInResponse
	SignInResponse struct {
		Result *apierror.APIError `json:"result"`
		Data   *SignInData        `json:"data"`
	}

	SignInData struct {
		IsSuccessfully bool   `json:"isSuccessfully"`
		Token          string `json:"token"`
	}
)

// create reservation models
type (
	CreateReservationRequest struct {
		Header
		IPAddress     string `json:"-"`
		UserId        int64  `json:"-"`
		Destination   string `json:"destination" validate:"required"`
		CheckInDate   string `json:"checkInDate" validate:"required"`
		CheckOutDate  string `json:"checkOutDate" validate:"required"`
		Accommodation string `json:"accommodation" validate:"required"`
		GuestCount    int    `json:"guestCount" validate:"required"`
	}

	CreateReservationResponse struct {
		Result *apierror.APIError     `json:"result"`
		Data   *CreateReservationData `json:"data"`
	}

	CreateReservationData struct {
		IsSuccessfully bool   `json:"isSuccessfully"`
		PNR            string `json:"pnr"`
	}
)

// update reservation models
type (
	UpdateReservationRequest struct {
		Header
		IPAddress     string `json:"-"`
		UserId        int64  `json:"-"`
		PNR           string `json:"pnr" validate:"required"`
		Destination   string `json:"destination" validate:"required"`
		CheckInDate   string `json:"checkInDate" validate:"required"`
		CheckOutDate  string `json:"checkOutDate" validate:"required"`
		Accommodation string `json:"accommodation" validate:"required"`
		GuestCount    int    `json:"guestCount" validate:"required"`
	}

	UpdateReservationResponse struct {
		Result *apierror.APIError     `json:"result"`
		Data   *UpdateReservationData `json:"data"`
	}

	UpdateReservationData struct {
		IsSuccessfully bool `json:"isSuccessfully"`
	}
)

// find reservation models
type (
	// FindReservationRequest defines the request structure for finding reservations.
	FindReservationRequest struct {
		Header
		IPAddress string `json:"-"`
		UserId    int64  `json:"-"`
		PNR       string `json:"-" query:"pnr" validate:"required"`
	}

	// FindReservationResponse defines the response structure for finding reservations.
	FindReservationResponse struct {
		Result *apierror.APIError   `json:"result"`
		Data   *FindReservationData `json:"data"`
	}

	// FindReservationData defines the details of a reservation.
	FindReservationData struct {
		PNR           string `json:"pnr"`
		Destination   string `json:"destination"`
		CheckInDate   string `json:"checkInDate"`
		CheckOutDate  string `json:"checkOutDate"`
		Accommodation string `json:"accommodation"`
		GuestCount    int    `json:"guestCount"`
		UserName      string `json:"userName"`
	}
)

// FindReservations
type (
	FindReservationsRequest struct {
		Header
		IPAddress string `json:"-"`
		UserId    int64  `json:"-"`
	}

	FindReservationsResponse struct {
		Result *apierror.APIError    `json:"result"`
		Data   *FindReservationsData `json:"data"`
	}

	FindReservationsData struct {
		Reservations   []FindReservationData `json:"reservations"`
		IsSuccessfully bool                  `json:"isSuccessfully"`
	}
)

// Localize method for FindReservationsResponse
func (f FindReservationsResponse) Localize(l *i18n.Localizer) interface{} {
	return f
}

// APIError method for FindReservationsResponse
func (f FindReservationsResponse) APIError() error {
	if f.Result == nil {
		return nil
	}

	return f.Result
}

// SetIPAddress method for FindReservationsRequest
func (f FindReservationsRequest) SetIPAddress(ipAddress string) {
	f.IPAddress = ipAddress
}

// Localize method for FindReservationResponse
func (f FindReservationResponse) Localize(l *i18n.Localizer) interface{} {
	return f
}

// APIError method for FindReservationResponse
func (f FindReservationResponse) APIError() error {
	if f.Result == nil {
		return nil
	}

	return f.Result
}

// SetIPAddress method for FindReservationRequest
func (f *FindReservationRequest) SetIPAddress(ipAddress string) {
	f.IPAddress = ipAddress
}

func (u UpdateReservationResponse) Localize(l *i18n.Localizer) interface{} {
	return u
}

func (u UpdateReservationResponse) APIError() error {
	if u.Result == nil {
		return nil
	}

	return u.Result
}

func (u UpdateReservationRequest) SetIPAddress(ipAddress string) {
	u.IPAddress = ipAddress
}

func (c CreateReservationResponse) Localize(l *i18n.Localizer) interface{} {
	return c
}

func (c CreateReservationResponse) APIError() error {
	if c.Result == nil {
		return nil
	}

	return c.Result
}

func (c CreateReservationRequest) SetIPAddress(ipAddress string) {
	c.IPAddress = ipAddress
}

func (p SignInResponse) Localize(_ *i18n.Localizer) interface{} {
	return p
}

func (p SignInResponse) APIError() error {
	if p.Result == nil {
		return nil
	}

	return p.Result
}

func (p SignInRequest) SetIPAddress(ipAddress string) {
	p.IPAddress = ipAddress
}

func (r HealthRequest) SetIPAddress(ipAddress string) {
	r.IPAddress = ipAddress
}

func (r HealthResponse) Localize(_ *i18n.Localizer) interface{} {
	return r
}

func (r HealthResponse) APIError() error {
	return nil
}
