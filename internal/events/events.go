package events

type OrderPlacedEvent struct {
    OrderID   string
    SellerID  string
    UserEmail string
}
