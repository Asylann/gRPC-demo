package models

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Size        int     `json:"size"`
	CategoryID  int     `json:"category_id"`
	ImageURL    string  `json:"imageURL"`
	SellerID    int     `json:"seller_id"`
}
