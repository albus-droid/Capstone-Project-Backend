package listing

// Service defines what our handlers expect.
type Service interface {
    Create(l *Listing) error
    GetByID(id string) (*Listing, error)
    ListBySeller(sellerID string) ([]Listing, error)
    ListAll() []Listing

    // new
    Update(id string, l Listing) error
    Delete(id string) error
}
