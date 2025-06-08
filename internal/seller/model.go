package seller

type Seller struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Phone    string `json:"phone"`
    Verified bool   `json:"verified"`
}
