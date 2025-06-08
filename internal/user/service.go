package user

// Service defines what we can do with users
type Service interface {
  Register(u User) error
  Authenticate(email, password string) (*User, error)
  GetByEmail(email string) (*User, error)
}
