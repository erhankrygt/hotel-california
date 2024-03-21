package authmiddlaware

import (
	"context"
	"errors"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/golang-jwt/jwt/v4"
	hotelcalifornia "hotel-california-backend"
	"hotel-california-backend/configs/envvars"
	apierror "hotel-california-backend/internal/api-error"
	"hotel-california-backend/internal/middleware"
	"strconv"
	"time"
)

// errors
var (
	errTokenHasExpired = errors.New("token has expired")
	errParsingUserId   = errors.New("error while parsing user ID")
)

// AuthMiddleware represents auth service middleware
type AuthMiddleware struct {
	l    log.Logger
	next hotelcalifornia.Service
	jw   envvars.JWTToken
}

// NewAuthMiddleware creates and returns auth middleware
func NewAuthMiddleware(l log.Logger, jw envvars.JWTToken) middleware.Middleware {
	return func(next hotelcalifornia.Service) hotelcalifornia.Service {
		return &AuthMiddleware{
			l:    l,
			next: next,
			jw:   jw,
		}
	}
}

// Health represents auth middleware's health method
func (m *AuthMiddleware) Health(ctx context.Context, req hotelcalifornia.HealthRequest) (res hotelcalifornia.HealthResponse) {
	return m.next.Health(ctx, req)
}

// SignIn represents auth middleware's sign in method
func (m *AuthMiddleware) SignIn(ctx context.Context, req hotelcalifornia.SignInRequest) hotelcalifornia.SignInResponse {
	return m.next.SignIn(ctx, req)
}

func (m *AuthMiddleware) CreateReservation(ctx context.Context, req hotelcalifornia.CreateReservationRequest) (res hotelcalifornia.CreateReservationResponse) {
	token := req.Token
	userId, err := m.isTokenValid(token)
	if err != nil {
		m.log(err, map[string]interface{}{
			"method": "CreateReservation",
			"action": "isTokenValid",
			"token":  token,
			"error":  err.Error(),
		})

		res.Result = apierror.DefaultUnauthorizedError
		return res
	}

	req.UserId = userId

	return m.next.CreateReservation(ctx, req)
}

func (m *AuthMiddleware) UpdateReservation(ctx context.Context, req hotelcalifornia.UpdateReservationRequest) (res hotelcalifornia.UpdateReservationResponse) {
	token := req.Token
	userId, err := m.isTokenValid(token)
	if err != nil {
		m.log(err, map[string]interface{}{
			"method": "UpdateReservation",
			"action": "isTokenValid",
			"token":  token,
			"error":  err.Error(),
		})

		res.Result = apierror.DefaultUnauthorizedError
		return res
	}

	req.UserId = userId

	return m.next.UpdateReservation(ctx, req)
}

func (m *AuthMiddleware) FindReservation(ctx context.Context, req hotelcalifornia.FindReservationRequest) (res hotelcalifornia.FindReservationResponse) {
	token := req.Token
	userId, err := m.isTokenValid(token)
	if err != nil {
		m.log(err, map[string]interface{}{
			"method": "FindReservation",
			"action": "isTokenValid",
			"token":  token,
			"error":  err.Error(),
		})

		res.Result = apierror.DefaultUnauthorizedError
		return res
	}

	req.UserId = userId

	return m.next.FindReservation(ctx, req)
}

func (m *AuthMiddleware) FindReservations(ctx context.Context, req hotelcalifornia.FindReservationsRequest) (res hotelcalifornia.FindReservationsResponse) {
	token := req.Token
	userId, err := m.isTokenValid(token)
	if err != nil {
		m.log(err, map[string]interface{}{
			"method": "FindReservations",
			"action": "isTokenValid",
			"token":  token,
			"error":  err.Error(),
		})

		res.Result = apierror.DefaultUnauthorizedError
		return res
	}

	req.UserId = userId

	return m.next.FindReservations(ctx, req)
}

func (m *AuthMiddleware) isTokenValid(tokenString string) (int64, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.jw.Secret), nil
	})

	// Check if there's an error in parsing or token is not valid
	if err != nil || !token.Valid {
		return 0, err
	}

	// Extract claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, err
	}

	// Check the expiration time of the token
	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
	if time.Now().After(expirationTime) {
		return 0, errTokenHasExpired
	}

	// Extract user ID from claims
	userid, err := strconv.ParseInt(claims["sub"].(string), 10, 64)
	if err != nil {
		return 0, errParsingUserId
	}

	// Return user ID and true indicating token is valid
	return userid, nil
}

func (s *AuthMiddleware) log(err error, additionalParams map[string]interface{}) {
	logParams := make([]interface{}, 0, 2+len(additionalParams)*2)

	for k, v := range additionalParams {
		logParams = append(logParams, k, v)
	}

	logParams = append(logParams, "error", err.Error())

	_ = level.Error(s.l).Log(logParams...)
}
