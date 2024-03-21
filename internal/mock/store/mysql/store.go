package mysqlstoretmock

import (
	"context"
	"github.com/stretchr/testify/mock"
	mysqlstore "hotel-california-backend/internal/store/mysql"
)

// compile-time proof of mongo store interface implementation
var _ mysqlstore.Store = (*Store)(nil)

// Store represents mock mongo store
type Store struct {
	mock.Mock
}

// NewStore returns mock mysql store
func NewStore() *Store {
	return &Store{}
}

func (s *Store) SignIn(ctx context.Context, username, password string) (usr *mysqlstore.User, err error) {
	args := s.Called(ctx, username, password)
	return args.Get(0).(*mysqlstore.User), args.Error(1)
}

func (s *Store) CreateReservation(ctx context.Context, res *mysqlstore.Reservation) error {
	args := s.Called(ctx, res)
	return args.Error(0)
}

func (s *Store) UpdateReservation(ctx context.Context, res *mysqlstore.Reservation) error {
	args := s.Called(ctx, res)
	return args.Error(0)
}

func (s *Store) FindReservation(ctx context.Context, pnr string, userID int64) (*mysqlstore.Reservation, error) {
	args := s.Called(ctx, pnr, userID)
	return args.Get(0).(*mysqlstore.Reservation), args.Error(1)
}

func (s *Store) FindReservations(ctx context.Context, userID int64) ([]*mysqlstore.Reservation, error) {
	args := s.Called(ctx, userID)
	return args.Get(0).([]*mysqlstore.Reservation), args.Error(1)
}

func (s *Store) Close() error {
	args := s.Called()
	return args.Error(0)
}
