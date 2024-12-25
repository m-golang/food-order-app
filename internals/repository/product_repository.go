package repository

import (
	"database/sql"

	"github.com/m-golang/food-order-app/internals/models"
)

type RepoModel struct {
	DB *sql.DB
}

// Retrieves a list of products for a specific menu
func (m *RepoModel) GetProducts(nameMenu string) ([]*models.Product, error) {
	query := `SELECT id, product_name, product_description, product_price, product_image FROM products WHERE menu_id = (SELECT id FROM menu WHERE name = ?)`

	rows, err := m.DB.Query(query, nameMenu)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := []*models.Product{}

	for rows.Next() {
		p := &models.Product{}

		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageName)
		if err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	return products, nil
}

// Retrieves a product by its ID
func (m *RepoModel) GetProductByID(id int) (*models.OrderProduct, error) {
	query := `SELECT id, product_name, product_price FROM products WHERE id = ?`
	orderProduct := &models.OrderProduct{}

	err := m.DB.QueryRow(query, id).Scan(&orderProduct.ID, &orderProduct.PName, &orderProduct.PPrice)
	if err != nil {
		return nil, err
	}

	return orderProduct, nil
}
