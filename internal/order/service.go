// internal/order/service.go
package order
type Service interface {
    Create(o Order) error
    GetByID(id string) (*Order, error)
    ListByUser(email string) ([]Order, error)

    Accept(id, callerEmail string) error   // ← 2 args
    Complete(id, callerEmail string) error // ← 2 args
}
