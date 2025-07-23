package user

type User struct {
  ID       string `json:"id" gorm:"primaryKey"`
  Name     string `json:"name"`
  Email    string `json:"email" gorm:"uniqueIndex;not null"`
  Password string `json:"password"` // Store hashed passwords
}