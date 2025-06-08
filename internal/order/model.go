package order

type Order struct {
    ID        string   `json:"id"`
    UserID    string   `json:"userId"`
    ListingIDs []string `json:"listingIds"`
    Total     float64  `json:"total"`
    CreatedAt int64    `json:"createdAt"` // Unix timestamp
}
