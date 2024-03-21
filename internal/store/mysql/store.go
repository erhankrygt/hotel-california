package mysqlstore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

type User struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement"`
	FirstName string    `gorm:"column:first_name"`
	LastName  string    `gorm:"column:last_name"`
	Username  string    `gorm:"column:username"`
	Password  string    `gorm:"column:password"`
	CreatedAt time.Time `gorm:"column:createdAt"`
	IsActive  bool      `gorm:"column:is_active"`
	IsDeleted bool      `gorm:"column:is_deleted"`
}

type Reservation struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement"`
	UserID        int64     `gorm:"column:user_id"`
	User          User      `gorm:"foreignKey:UserID"`
	PNR           string    `gorm:"column:pnr"`
	Destination   string    `gorm:"column:destination"`
	CheckInDate   time.Time `gorm:"column:check_in_date"`
	CheckOutDate  time.Time `gorm:"column:check_out_date"`
	Accommodation string    `gorm:"column:accommodation"`
	GuestCount    int       `gorm:"column:guest_count"`
	CreatedAt     time.Time `gorm:"column:createdAt"`
	IsActive      bool      `gorm:"column:is_active"`
	IsDeleted     bool      `gorm:"column:is_deleted"`
}

type Store interface {
	SignIn(ctx context.Context, username, password string) (usr *User, err error)
	CreateReservation(ctx context.Context, res *Reservation) error
	UpdateReservation(ctx context.Context, res *Reservation) error
	FindReservation(ctx context.Context, pnr string, userID int64) (*Reservation, error)
	FindReservations(ctx context.Context, userID int64) ([]*Reservation, error)
	Close() error
}

// store implements Store interface
type store struct {
	db   *gorm.DB
	opts Options
}

// compile-time proof of interface implementation
var _ Store = (*store)(nil)

type Options struct {
	URI               string
	Database          string
	UserName          string
	Password          string
	Port              string
	SSLMode           string
	ConnectTimeout    int
	PingTimeout       time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	DisconnectTimeout time.Duration
}

// NewStore creates and returns collect store
func NewStore(opts Options) (Store, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", opts.UserName, opts.Password, opts.URI, opts.Port, opts.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&User{}, Reservation{})
	if err != nil {
		return nil, err
	}

	cli := store{
		db:   db,
		opts: opts,
	}

	return &cli, nil
}

func (c *store) SignIn(ctx context.Context, username, password string) (usr *User, err error) {
	query := "username= ? and password = ? and is_active = ? and is_deleted = ?"
	err = c.db.WithContext(ctx).Model(&User{}).Where(query, username, password, true, false).Find(&usr).Error
	if err != nil {
		return nil, err
	}

	if usr.ID == 0 {
		return nil, errors.New("user not found")
	}

	return usr, nil
}

func (s *store) CreateReservation(ctx context.Context, res *Reservation) error {
	err := s.db.WithContext(ctx).Create(res).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *store) UpdateReservation(ctx context.Context, res *Reservation) error {
	tx := s.db.Begin()

	if tx.Error != nil {
		return tx.Error
	}

	var reservation Reservation
	if err := tx.WithContext(ctx).Where("pnr = ?", res.PNR).First(&reservation).Error; err != nil {
		tx.Rollback()
		return err
	}

	if res.UserID != reservation.UserID {
		return errors.New("the reservation you are trying to update does not belong to this user")
	}

	if time.Now().After(reservation.CheckInDate) {
		return errors.New("the check-in date of a past reservation cannot be changed")
	}

	reservation.Destination = res.Destination
	reservation.CheckInDate = res.CheckInDate
	reservation.CheckOutDate = res.CheckOutDate
	reservation.Accommodation = res.Accommodation
	reservation.GuestCount = res.GuestCount

	if err := tx.WithContext(ctx).Save(&reservation).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (s *store) FindReservation(ctx context.Context, pnr string, userID int64) (*Reservation, error) {
	var reservation Reservation
	if err := s.db.WithContext(ctx).Where("pnr = ? AND user_id = ?", pnr, userID).Preload("User").First(&reservation).Error; err != nil {
		return nil, err
	}

	return &reservation, nil
}

func (s *store) FindReservations(ctx context.Context, userID int64) ([]*Reservation, error) {
	var reservations []*Reservation

	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).
		Preload("User").
		Find(&reservations).
		Error; err != nil {
		return nil, err
	}

	return reservations, nil
}

// Close returns database close
func (c *store) Close() error {
	db, err := c.db.DB()
	if err != nil {
		return err
	}

	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	return nil
}
