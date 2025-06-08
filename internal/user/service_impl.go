package user

import (
  "errors"
  "golang.org/x/crypto/bcrypt"
)

type inMemoryService struct {
  users map[string]User
}

func NewInMemoryService() Service {
  return &inMemoryService{users: make(map[string]User)}
}

func (s *inMemoryService) Register(u User) error {
  if _, ok := s.users[u.Email]; ok {
    return errors.New("user already exists")
  }
  // hash password
  h, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
  if err != nil {
    return err
  }
  u.Password = string(h)
  s.users[u.Email] = u
  return nil
}

func (s *inMemoryService) Authenticate(email, pw string) (*User, error) {
  u, err := s.GetByEmail(email)
  if err != nil {
    return nil, err
  }
  if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pw)); err != nil {
    return nil, errors.New("invalid credentials")
  }
  return u, nil
}

func (s *inMemoryService) GetByEmail(email string) (*User, error) {
  u, ok := s.users[email]
  if !ok {
    return nil, errors.New("user not found")
  }
  return &u, nil
}
