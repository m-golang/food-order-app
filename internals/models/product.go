package models

// Product represents a product in the menu
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"product_name"`
	Description string  `json:"product_description"`
	Price       float64 `json:"product_price"`
	ImageName   string  `json:"product_image"`
}

// OrderProduct represents a product purchased in an order
type OrderProduct struct {
	ID     int     `json:"id"`
	PName  string  `json:"product_name"`
	PPrice float64 `json:"product_price"`
}
