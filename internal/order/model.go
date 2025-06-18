package order

type Order struct {
    ID        string `json:"id"`
    UserID    string `json:"userId"`
    UserEmail string `json:"userEmail"`
    SellerID  string `json:"sellerId"`
    ItemID    string `json:"itemId"`
}
