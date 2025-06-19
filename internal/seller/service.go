package seller

// Service defines seller operations
type Service interface {
    Register(s Seller) error
    GetByID(id string) (*Seller, error)
    ListAll() []Seller
}
