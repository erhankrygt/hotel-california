package endpoints

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	hotelcalifornia "hotel-california-backend"
)

// Endpoints represents service endpoints
type Endpoints struct {
	HealthEndpoint            endpoint.Endpoint
	SignInEndpoint            endpoint.Endpoint
	CreateReservationEndpoint endpoint.Endpoint
	UpdateReservationEndpoint endpoint.Endpoint
	FindReservationEndpoint   endpoint.Endpoint
	FindReservationsEndpoint  endpoint.Endpoint
}

// MakeEndpoints makes and returns endpoints
func MakeEndpoints(s hotelcalifornia.Service) Endpoints {
	return Endpoints{
		HealthEndpoint:            MakeHealthEndpoint(s),
		SignInEndpoint:            MakeSignInEndpoint(s),
		CreateReservationEndpoint: MakeCreateReservationEndpoint(s),
		UpdateReservationEndpoint: MakeUpdateReservationEndpoint(s),
		FindReservationEndpoint:   MakeFindReservationEndpoint(s),
		FindReservationsEndpoint:  MakeFindReservationsEndpoint(s),
	}
}

// MakeHealthEndpoint makes and returns health endpoint
func MakeHealthEndpoint(s hotelcalifornia.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*hotelcalifornia.HealthRequest)

		res := s.Health(ctx, *req)

		return res, nil
	}
}

// MakeSignInEndpoint makes and returns sign in endpoint
func MakeSignInEndpoint(s hotelcalifornia.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*hotelcalifornia.SignInRequest)

		res := s.SignIn(ctx, *req)

		return res, nil
	}
}

// MakeCreateReservationEndpoint makes and returns create reservation endpoint
func MakeCreateReservationEndpoint(s hotelcalifornia.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*hotelcalifornia.CreateReservationRequest)

		res := s.CreateReservation(ctx, *req)

		return res, nil
	}
}

// MakeUpdateReservationEndpoint makes and returns update reservation endpoint
func MakeUpdateReservationEndpoint(s hotelcalifornia.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*hotelcalifornia.UpdateReservationRequest)

		res := s.UpdateReservation(ctx, *req)

		return res, nil
	}
}

// MakeFindReservationEndpoint makes and returns find reservation endpoint
func MakeFindReservationEndpoint(s hotelcalifornia.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*hotelcalifornia.FindReservationRequest)

		res := s.FindReservation(ctx, *req)

		return res, nil
	}
}

// MakeFindReservationsEndpoint makes and returns find reservations endpoint
func MakeFindReservationsEndpoint(s hotelcalifornia.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*hotelcalifornia.FindReservationsRequest)

		res := s.FindReservations(ctx, *req)

		return res, nil
	}
}
