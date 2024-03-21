package httptransport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/iris-contrib/schema"
	hotelcalifornia "hotel-california-backend"
	apierror "hotel-california-backend/internal/api-error"
	"hotel-california-backend/internal/endpoints"
	"hotel-california-backend/internal/localization"
	"hotel-california-backend/internal/transport"
	"net/http"
	"reflect"
)

// endpoint names
const (
	health            = "Health"
	signIn            = "SignIn"
	createReservation = "CreateReservation"
	updateReservation = "UpdateReservation"
	findReservation   = "FindReservation"
	findReservations  = "FindReservations"
)

// decoder tags
const (
	headerTag = "header"
	queryTag  = "query"
)

const invalidResponseError = "invalid response"

// MakeHTTPHandler makes and returns http handler
func MakeHTTPHandler(l log.Logger, s hotelcalifornia.Service) http.Handler {
	es := endpoints.MakeEndpoints(s)

	r := mux.NewRouter()

	// GET /health
	r.Methods(http.MethodGet).Path("/health").Handler(
		makeHealthHandler(es.HealthEndpoint, makeDefaultServerOptions(l, health)),
	)

	// hotel california router
	router := r.PathPrefix("/v1/").Subrouter()

	// POST /account/sign-in
	router.Methods(http.MethodPost).Path("/account/sign-in").Handler(
		makeSignInHandler(es.SignInEndpoint, makeDefaultServerOptions(l, signIn)),
	)

	// POST /reservation/new
	router.Methods(http.MethodPost).Path("/reservation/new").Handler(
		makeCreateReservationHandler(es.CreateReservationEndpoint, makeDefaultServerOptions(l, createReservation)),
	)

	// POST /reservation/update
	router.Methods(http.MethodPost).Path("/reservation/update").Handler(
		makeUpdateReservationHandler(es.UpdateReservationEndpoint, makeDefaultServerOptions(l, updateReservation)),
	)

	// GET /reservation
	router.Methods(http.MethodGet).Path("/reservation").Handler(
		makeFindReservationHandler(es.FindReservationEndpoint, makeDefaultServerOptions(l, findReservation)),
	)

	// GET /reservation
	router.Methods(http.MethodGet).Path("/reservations").Handler(
		makeFindReservationsHandler(es.FindReservationsEndpoint, makeDefaultServerOptions(l, findReservations)),
	)

	return r
}

func makeHealthHandler(e endpoint.Endpoint, serverOptions []kithttp.ServerOption) http.Handler {
	h := kithttp.NewServer(e, makeDecoder(hotelcalifornia.HealthRequest{}), encoder, serverOptions...)

	return h
}

func makeSignInHandler(e endpoint.Endpoint, serverOptions []kithttp.ServerOption) http.Handler {
	h := kithttp.NewServer(e, makeDecoder(hotelcalifornia.SignInRequest{}), encoder, serverOptions...)

	return h
}

func makeCreateReservationHandler(e endpoint.Endpoint, serverOptions []kithttp.ServerOption) http.Handler {
	h := kithttp.NewServer(e, makeDecoder(hotelcalifornia.CreateReservationRequest{}), encoder, serverOptions...)

	return h
}

func makeUpdateReservationHandler(e endpoint.Endpoint, serverOptions []kithttp.ServerOption) http.Handler {
	h := kithttp.NewServer(e, makeDecoder(hotelcalifornia.UpdateReservationRequest{}), encoder, serverOptions...)

	return h
}

func makeFindReservationHandler(e endpoint.Endpoint, serverOptions []kithttp.ServerOption) http.Handler {
	h := kithttp.NewServer(e, makeDecoder(hotelcalifornia.FindReservationRequest{}), encoder, serverOptions...)

	return h
}

func makeFindReservationsHandler(e endpoint.Endpoint, serverOptions []kithttp.ServerOption) http.Handler {
	h := kithttp.NewServer(e, makeDecoder(hotelcalifornia.FindReservationsRequest{}), encoder, serverOptions...)

	return h
}

func makeDefaultServerOptions(l log.Logger, endpointName string) []kithttp.ServerOption {
	return []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(errorEncoder),
		kithttp.ServerErrorHandler(transport.NewErrorHandler(l, endpointName)),
		kithttp.ServerBefore(localization.AddLocalizerToContext),
	}
}

func makeDecoder(emptyReq interface{}) kithttp.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		req := reflect.New(reflect.TypeOf(emptyReq)).Interface()

		if err := newHeaderDecoder().Decode(req, r.Header); err != nil {
			return nil, fmt.Errorf("decoding request header failed, %s", err.Error())
		}

		if err := newQueryDecoder().Decode(req, r.URL.Query()); err != nil {
			return nil, fmt.Errorf("decoding request query failed, %s", err.Error())
		}

		if requestHasBody(r) {
			if err := json.NewDecoder(r.Body).Decode(req); err != nil {
				return nil, fmt.Errorf("decoding request body failed, %s", err.Error())
			}
		}

		if err := validate(req); err != nil {
			apiError := apierror.NewValidationError(err.Error(), "")
			apiError.BaseError = err
			return nil, apiError
		}

		return req, nil
	}
}

func newHeaderDecoder() *schema.Decoder {
	return newDecoder(headerTag)
}

func newQueryDecoder() *schema.Decoder {
	return newDecoder(queryTag)
}

func newDecoder(tag string) *schema.Decoder {
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	if tag != "" {
		decoder.SetAliasTag(tag)
	}

	return decoder
}

func requestHasBody(r *http.Request) bool {
	return r.Body != http.NoBody
}

func validate(req interface{}) error {
	errs := validator.New().Struct(req)
	if errs == nil {
		return nil
	}

	firstErr := errs.(validator.ValidationErrors)[0]

	return errors.New("validation failed, tag: " + firstErr.Tag() + ", field: " + firstErr.Field())
}

func encoder(ctx context.Context, rw http.ResponseWriter, response interface{}) error {
	r, ok := response.(hotelcalifornia.Response)
	if !ok {
		return errors.New(invalidResponseError)
	}

	if r.APIError() != nil {
		errorEncoder(ctx, r.APIError(), rw)

		return nil
	}

	l := localization.GetLocalizerFromContext(ctx)

	lr := r.Localize(l)

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	return json.NewEncoder(rw).Encode(lr)
}

func errorEncoder(ctx context.Context, err error, rw http.ResponseWriter) {
	var apiErr *apierror.APIError

	ok := errors.As(err, &apiErr)
	if !ok {
		apiErr = apierror.DefaultInternalServerError
	}

	l := localization.GetLocalizerFromContext(ctx)

	apiErr.Localize(l)

	er := errorResponse{
		Data:   nil,
		Result: apiErr,
	}

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(apiErr.StatusCode)

	_ = json.NewEncoder(rw).Encode(er)
}

type errorResponse struct {
	Data   interface{}        `json:"data"`
	Result *apierror.APIError `json:"result"`
}
