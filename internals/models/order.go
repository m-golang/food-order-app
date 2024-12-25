package models

// Orders represents an order with a list of associated products (OrderItems)
type Orders struct {
	OrderID         int
	TotalAmount     float64
	Status          string
	OrderProducts []*OrderItem
}

// OrderItem represents an individual product in an order
type OrderItem struct {
	ProductName string
	Quantity    int
}
