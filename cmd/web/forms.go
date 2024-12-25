package main

// SignupForm represents the structure of the form data for user registration.
type SignupForm struct {
	FullName    string `form:"full_name" json:"full_name" binding:"required"`
	PhoneNumber string `form:"phone_number" json:"phone_number" binding:"required"`
	Password    string `form:"password" json:"password" binding:"required"`
}

// LoginForm represents the structure of the form data for user login.
type LoginForm struct {
	PhoneNumber string `form:"phone_number" json:"phone_number" binding:"required"`
	Password    string `form:"password" json:"password" binding:"required"`
}

// UserUpdateForm is used for updating the user's full name.
type UserUpdateForm struct {
	FullName string `form:"full_name" json:"full_name" binding:"required"`
}


// BasketItem represents a single item in the user's shopping basket.
type BasketItem struct {
	ID       int `form:"id" json:"id" binding:"required"`
	Quantity int `form:"quantity" json:"quantity" binding:"required"`
}

// BasketOrder represents the structure for placing an order, including multiple basket items.
type BasketOrder struct {
	OrderProducts   []BasketItem `form:"order_products" json:"order_products" binding:"required"`
	DeliveryAddress string       `form:"delivery_address" json:"delivery_address" binding:"required"`
}

// OrderProducts represents the relationship between orders and the products that belong to them.
type OrderProducts struct {
	OrderID   int
	ProductID int
	Quantity  int
}
