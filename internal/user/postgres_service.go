package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type postgresService struct {
	db *gorm.DB
}

func NewPostgresService(db *gorm.DB) Service {
	return &postgresService{db: db}
}

func (s *postgresService) Register(u User) error {
	// Check if user exists
	var existing User
	if err := s.db.Where("email = ?", u.Email).First(&existing).Error; err == nil {
		return errors.New("user already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Hash the password
	h, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(h)

	// Save to DB
	return s.db.Create(&u).Error
}

func (s *postgresService) Authenticate(email, pw string) (*User, error) {
	var u User
	if err := s.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, errors.New("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pw)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return &u, nil
}

func (s *postgresService) GetByEmail(email string) (*User, error) {
	var u User
	if err := s.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, errors.New("user not found")
	}
	return &u, nil
}
