package listing

// Service defines listing operations
type Service interface {
    Create(l Listing) error
    GetByID(id string) (*Listing, error)
    ListBySeller(sellerID string) ([]Listing, error)
    ListAll() []Listing
}
