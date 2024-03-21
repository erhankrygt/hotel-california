package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/golang-jwt/jwt/v4"
	hotelcalifornia "hotel-california-backend"
	"hotel-california-backend/configs/envvars"
	apierror "hotel-california-backend/internal/api-error"
	mysqlstore "hotel-california-backend/internal/store/mysql"
	"math/rand"
	"strconv"
	"time"
)

var (
	AccommodationTypes = []string{
		"beach",
		"city",
		"mountain",
	}
	dateLayout = "2006-01-02"
	pnrLen     = 8
	charset    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// errors
var (
	ErrNotValidDestination = errors.New("destination not valid")
	ErrFailedParseDate     = errors.New("failed to parse date")
)

// compile-time proof of service interface implementation
var _ hotelcalifornia.Service = (*Service)(nil)

// Service represents service
type Service struct {
	environment string
	l           log.Logger
	ms          mysqlstore.Store
	jw          envvars.JWTToken
}

// NewService creates and returns service
func NewService(environment string, l log.Logger, ms mysqlstore.Store, jw envvars.JWTToken) hotelcalifornia.Service {
	return &Service{
		environment: environment,
		l:           l,
		ms:          ms,
		jw:          jw,
	}
}

// Health represents service's health method
func (s *Service) Health(_ context.Context, _ hotelcalifornia.HealthRequest) hotelcalifornia.HealthResponse {
	res := hotelcalifornia.HealthResponse{
		Ping: "Pong",
	}

	return res
}

// SignIn represents service's sign in method
func (s *Service) SignIn(ctx context.Context, req hotelcalifornia.SignInRequest) hotelcalifornia.SignInResponse {
	res := hotelcalifornia.SignInResponse{}

	userName := req.UserName

	hash := md5.Sum([]byte(req.Password))
	password := hex.EncodeToString(hash[:])

	usr, err := s.ms.SignIn(ctx, userName, password)
	if err != nil {
		s.log(err, map[string]interface{}{
			"method": "SignIn",
			"action": "Mysql SignIn",
			"error":  err.Error(),
		})

		res.Result = apierror.NewBadRequestError(err)
		res.Result.BaseError = err

		return res
	}

	token, err := s.createToken(usr.ID)
	if err != nil {
		s.log(err, map[string]interface{}{
			"method": "SignIn",
			"action": "CreateToken",
			"error":  err.Error(),
		})

		res.Result = apierror.NewBadRequestError(err)
		res.Result.BaseError = err

		return res
	}

	res.Data = &hotelcalifornia.SignInData{
		IsSuccessfully: true,
		Token:          token,
	}

	return res
}

// CreateReservation represents service's create reservation method
func (s *Service) CreateReservation(ctx context.Context, req hotelcalifornia.CreateReservationRequest) hotelcalifornia.CreateReservationResponse {
	res := hotelcalifornia.CreateReservationResponse{}
	userId := req.UserId

	va := s.validAccommodation(req.Accommodation)
	if va == false {
		res.Result = apierror.NewBadRequestError(ErrNotValidDestination)
		res.Result.BaseError = ErrNotValidDestination
		return res
	}

	checkInDate, err := s.parseDate(req.CheckInDate, dateLayout)
	if err != nil {
		res.Result = apierror.NewBadRequestError(err)
		res.Result.BaseError = err
		return res
	}

	checkOutDate, err := s.parseDate(req.CheckOutDate, dateLayout)
	if err != nil {
		res.Result = apierror.NewBadRequestError(err)
		res.Result.BaseError = err
		return res
	}

	if checkInDate.After(checkOutDate) {
		res.Result = apierror.CouldNotCheckInGreaterThanCheckout
		return res
	}

	pnr := s.generatePNR(pnrLen)

	rev := mysqlstore.Reservation{
		UserID:        userId,
		PNR:           pnr,
		Destination:   req.Destination,
		CheckInDate:   checkInDate,
		CheckOutDate:  checkOutDate,
		CreatedAt:     time.Now(),
		Accommodation: req.Accommodation,
		GuestCount:    req.GuestCount,
		IsActive:      true,
		IsDeleted:     false,
	}

	err = s.ms.CreateReservation(ctx, &rev)
	if err != nil {
		res.Result = apierror.CouldNotCreateReservation
		res.Result.BaseError = err
		return res
	}

	res.Data = &hotelcalifornia.CreateReservationData{
		IsSuccessfully: true,
		PNR:            pnr,
	}

	return res
}

// UpdateReservation represents service's update reservation method
func (s *Service) UpdateReservation(ctx context.Context, req hotelcalifornia.UpdateReservationRequest) hotelcalifornia.UpdateReservationResponse {
	res := hotelcalifornia.UpdateReservationResponse{}

	userId := req.UserId
	pnr := req.PNR

	va := s.validAccommodation(req.Accommodation)
	if va == false {
		res.Result = apierror.NewBadRequestError(ErrNotValidDestination)
		res.Result.BaseError = ErrNotValidDestination
		return res
	}

	checkInDate, err := s.parseDate(req.CheckInDate, dateLayout)
	if err != nil {
		res.Result = apierror.NewBadRequestError(err)
		res.Result.BaseError = err
		return res
	}

	checkOutDate, err := s.parseDate(req.CheckOutDate, dateLayout)
	if err != nil {
		res.Result = apierror.NewBadRequestError(err)
		res.Result.BaseError = err
		return res
	}

	if time.Now().After(checkInDate) {
		res.Result = apierror.CouldNotChangeReservationCheckInDate
		return res
	}

	if checkInDate.After(checkOutDate) {
		res.Result = apierror.CouldNotCheckInGreaterThanCheckout
		return res
	}

	rev := mysqlstore.Reservation{
		PNR:           pnr,
		UserID:        userId,
		Destination:   req.Destination,
		CheckInDate:   checkInDate,
		CheckOutDate:  checkOutDate,
		Accommodation: req.Accommodation,
		GuestCount:    req.GuestCount,
	}

	err = s.ms.UpdateReservation(ctx, &rev)
	if err != nil {
		res.Result = apierror.NewBadRequestError(err)
		res.Result.BaseError = err
		return res
	}

	res.Data = &hotelcalifornia.UpdateReservationData{
		IsSuccessfully: true,
	}

	return res
}

// FindReservation represents service's find reservation method
func (s *Service) FindReservation(ctx context.Context, req hotelcalifornia.FindReservationRequest) hotelcalifornia.FindReservationResponse {
	res := hotelcalifornia.FindReservationResponse{}

	reservation, err := s.ms.FindReservation(ctx, req.PNR, req.UserId)
	if err != nil {
		res.Result = apierror.NewBadRequestError(err)
		res.Result.BaseError = err
		return res
	}

	res.Data = &hotelcalifornia.FindReservationData{
		PNR:           reservation.PNR,
		Destination:   reservation.Destination,
		CheckInDate:   reservation.CheckInDate.Format(dateLayout),
		CheckOutDate:  reservation.CheckOutDate.Format(dateLayout),
		Accommodation: reservation.Accommodation,
		GuestCount:    reservation.GuestCount,
		UserName:      fmt.Sprintf("%s %s", reservation.User.FirstName, reservation.User.LastName),
	}

	return res
}

// FindReservations represents service's find reservations method
func (s *Service) FindReservations(ctx context.Context, req hotelcalifornia.FindReservationsRequest) hotelcalifornia.FindReservationsResponse {
	res := hotelcalifornia.FindReservationsResponse{}

	reservations, err := s.ms.FindReservations(ctx, req.UserId)
	if err != nil {
		res.Result = apierror.NewBadRequestError(err)
		res.Result.BaseError = err
		return res
	}

	var rvs []hotelcalifornia.FindReservationData

	for _, reservation := range reservations {

		d := hotelcalifornia.FindReservationData{
			PNR:           reservation.PNR,
			Destination:   reservation.Destination,
			CheckInDate:   reservation.CheckInDate.Format(dateLayout),
			CheckOutDate:  reservation.CheckOutDate.Format(dateLayout),
			Accommodation: reservation.Accommodation,
			GuestCount:    reservation.GuestCount,
			UserName:      fmt.Sprintf("%s %s", reservation.User.FirstName, reservation.User.LastName),
		}

		rvs = append(rvs, d)
	}

	res.Data = &hotelcalifornia.FindReservationsData{
		Reservations:   rvs,
		IsSuccessfully: true,
	}

	return res
}

func (s *Service) createToken(userid int64) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		Subject:   strconv.FormatInt(userid, 10),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte(s.jw.Secret)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *Service) validAccommodation(ac string) bool {
	for _, a := range AccommodationTypes {
		if a == ac {
			return true
		}
	}

	return false
}

// It has nothing to do with the real scenario.
func (s *Service) generatePNR(length int) string {

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func (s *Service) parseDate(dateString, format string) (time.Time, error) {
	parsedTime, err := time.Parse(format, dateString)
	if err != nil {
		return time.Time{}, ErrFailedParseDate
	}

	return parsedTime, nil
}

func (s *Service) log(err error, additionalParams map[string]interface{}) {
	logParams := make([]interface{}, 0, 2+len(additionalParams)*2)

	for k, v := range additionalParams {
		logParams = append(logParams, k, v)
	}

	logParams = append(logParams, "error", err.Error())

	_ = level.Error(s.l).Log(logParams...)
}
