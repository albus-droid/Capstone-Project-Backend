package order

// Service defines order operations
type Service interface {
    Create(o Order) error
    GetByID(id string) (*Order, error)
    ListByUser(userID string) ([]Order, error)
}
