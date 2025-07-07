package seller

type Seller struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"` // accept on input, omitted when outputting
    Phone    string `json:"phone"`
    Verified bool   `json:"verified"`
}
