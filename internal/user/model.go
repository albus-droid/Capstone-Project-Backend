package user

// the JSON tags let Gin bind & render JSON automatically
type User struct {
  ID       string `json:"id"`
  Name     string `json:"name"`
  Email    string `json:"email"`
  Password string `json:"password"`   // we’ll hash it, never return it
}
