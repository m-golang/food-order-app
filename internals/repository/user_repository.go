package repository

import (
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/m-golang/food-order-app/internals/models"
	"golang.org/x/crypto/bcrypt"
)

// Inserts a new user into the database with hashed password
func (m *RepoModel) InsertUsers(fullName, phone_number, password string) error {
	query := `INSERT INTO users (full_name, phone_number,password_hash) VALUES (?, ?, ?)`

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = m.DB.Exec(query, fullName, phone_number, hashedPassword)
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 {
				return ErrDuplicatePhoneNumber
			}
		}
		return err
	}

	return nil
}

// Authenticates a user by validating the phone number and password
func (m *RepoModel) Authenticate(phoneNumber, password string) (int, error) {
	query := `SELECT id, password_hash FROM users WHERE phone_number = ?`
	var id int
	var password_hash string

	err := m.DB.QueryRow(query, phoneNumber).Scan(&id, &password_hash)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(password_hash), []byte(password))
	if err != nil {
		return 0, ErrInvalidCredentials
	}

	return id, nil
}

// Retrieves user info by user ID
func (m *RepoModel) GetUserInfoByID(userID int) (*models.UserInfo, error) {
	query := `SELECT full_name, phone_number FROM users WHERE id = ?`

	userInfo := &models.UserInfo{}

	err := m.DB.QueryRow(query, userID).Scan(&userInfo.FullName, &userInfo.PhoneNumber)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

// Updates the user's full name by their user ID
func (m *RepoModel) UpdateUserNameByID(userID int, userName string) error {
	query := `UPDATE users SET full_name = ? WHERE id = ?`

	_, err := m.DB.Exec(query, userName, userID)
	if err != nil {
		return err
	}
	return nil
}
