package listing

type Listing struct {
    ID          string  `json:"id"`
    SellerID    string  `json:"sellerId"`
    Title       string  `json:"title"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    Available   bool    `json:"available"`
}
