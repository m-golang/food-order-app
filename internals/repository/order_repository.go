package repository

import (
	"github.com/m-golang/food-order-app/internals/models"
)

// Creates a new order and returns its ID
func (m *RepoModel) CreateNewOrder(userID int, totalAmount float64, deliveryAddress string) (int, error) {
	query := `INSERT INTO orders (user_id, total_amount, delivery_address) VALUES (?, ?, ?)`

	result, err := m.DB.Exec(query, userID, totalAmount, deliveryAddress)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Creates new order items (products) associated with the given order
func (m *RepoModel) CreateNewOrderProducts(orderID, productID, quantity int) error {
	query := `INSERT INTO order_products (order_id, product_id, quantity) VALUES (?, ?, ?)`

	_, err := m.DB.Exec(query, orderID, productID, quantity)
	if err != nil {
		return err
	}

	return nil
}

// Retrieves all orders for a given user with associated products
func (m *RepoModel) GetOrdersWithItems(userID int) ([]*models.Orders, error) {
	// Query to get all orders for the user and join with order_products and products
	query := `SELECT o.id AS order_id, o.total_amount, o.status, p.product_name, op.quantity FROM orders o JOIN order_products op ON o.id = op.order_id JOIN products p ON op.product_id = p.id WHERE o.user_id = ?`

	rows, err := m.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	
	var orders []*models.Orders
	ordersMap := make(map[int]*models.Orders)

	// Iterate over the result set
	for rows.Next() {
		var (
			orderID             int
			totalAmount         float64
			status, productName string
			quantity            int
		)

		err := rows.Scan(&orderID, &totalAmount, &status, &productName, &quantity)
		if err != nil {
			return nil, err
		}

		// If order doesn't exist in the map, create a new one
		order, exists := ordersMap[orderID]
		if !exists {
			order = &models.Orders{
				OrderID:       orderID,
				TotalAmount:   totalAmount,
				Status:        status,
				OrderProducts: []*models.OrderItem{},
			}
			ordersMap[orderID] = order
		}

		// Append the order item (product) to the OrderProducts slice
		order.OrderProducts = append(order.OrderProducts, &models.OrderItem{
			ProductName: productName,
			Quantity:    quantity,
		})
	}

	// Convert the map to a slice for returning
	for _, order := range ordersMap {
		orders = append(orders, order)
	}

	return orders, nil
}
