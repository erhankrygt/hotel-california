package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
	hotelcalifornia "hotel-california-backend"
	"hotel-california-backend/configs/envvars"
	apierror "hotel-california-backend/internal/api-error"
	mysqlstoretmock "hotel-california-backend/internal/mock/store/mysql"
	mysqlstore "hotel-california-backend/internal/store/mysql"
	"os"
	"testing"
	"time"
)

func TestFindReservation_Success(t *testing.T) {
	// Context
	ctx := context.Background()

	// Log
	logger := log.NewLogfmtLogger(os.Stdout)

	// MySQL Mock
	ms := mysqlstoretmock.NewStore()

	pnr := "xyz123"
	var userId int64 = 1

	reservation := &mysqlstore.Reservation{
		CheckInDate:   time.Now(),
		CheckOutDate:  time.Now().AddDate(0, 0, 1),
		Accommodation: "mountain",
		Destination:   "Istanbul",
		GuestCount:    1,
		PNR:           pnr,
		UserID:        userId,
		User: mysqlstore.User{
			FirstName: "John",
			LastName:  "Doe",
		},
	}

	ms.On("FindReservation", ctx, pnr, userId).Return(reservation, nil)

	service := NewService("dev", logger, ms, envvars.JWTToken{})

	req := hotelcalifornia.FindReservationRequest{
		IPAddress: "127.192.1.1",
		PNR:       pnr,
		UserId:    userId,
	}

	expectedResponse := hotelcalifornia.FindReservationResponse{
		Data: &hotelcalifornia.FindReservationData{
			PNR:           reservation.PNR,
			Destination:   reservation.Destination,
			CheckInDate:   reservation.CheckInDate.Format(dateLayout),
			CheckOutDate:  reservation.CheckOutDate.Format(dateLayout),
			Accommodation: reservation.Accommodation,
			GuestCount:    reservation.GuestCount,
			UserName:      fmt.Sprintf("%s %s", reservation.User.FirstName, reservation.User.LastName),
		},
		Result: nil,
	}

	response := service.FindReservation(ctx, req)
	assert.Equal(t, expectedResponse, response)
}

func TestFindReservation_Error(t *testing.T) {
	// Context
	ctx := context.Background()

	// Log
	logger := log.NewLogfmtLogger(os.Stdout)

	// MySQL Mock
	ms := mysqlstoretmock.NewStore()

	pnr := "xyz123"
	var userId int64 = 1

	var reservation *mysqlstore.Reservation
	err := errors.New("record not found")

	ms.On("FindReservation", ctx, pnr, userId).Return(reservation, err)

	service := NewService("dev", logger, ms, envvars.JWTToken{})

	req := hotelcalifornia.FindReservationRequest{
		IPAddress: "127.192.1.1",
		PNR:       pnr,
		UserId:    userId,
	}

	expectedResponse := hotelcalifornia.FindReservationResponse{
		Data:   nil,
		Result: apierror.NewBadRequestError(err),
	}

	expectedResponse.Result.BaseError = err

	response := service.FindReservation(ctx, req)
	assert.Equal(t, expectedResponse, response)
}
