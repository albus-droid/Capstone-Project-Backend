package order

type Order struct {
	ID         string   `json:"id"`
	UserID     string   `json:"userId"`
	UserEmail string    `json:"user_email"` // 👈 new – owner’s e-mail
	SellerID   string   `json:"sellerId"` // NEW — who is fulfilling this order
	ListingIDs []string `json:"listingIds"`
	Total      float64  `json:"total"`
	CreatedAt  int64    `json:"createdAt"` // Unix timestamp
	Status     string   `json:"status"`    // NEW — e.g., "pending", "accepted", "rejected"
}