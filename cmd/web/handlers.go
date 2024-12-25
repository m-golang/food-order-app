package main

import (
	"errors"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/m-golang/food-order-app/internals/repository"
)

// GetBurgers retrieves all burger products from the repository and renders the view.
func (env *env) GetBurgers(c *gin.Context) {
	products, err := env.repo.GetProducts("burgers")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "base.html", gin.H{
			"error": http.StatusInternalServerError,
		})
		return
	}
	
	c.HTML(http.StatusOK, "base.html", gin.H{
		"title":    "Burgers",
		"Products": products,
	})
}

// GetFishes retrieves all fish products from the repository and renders the view.
func (env *env) GetFishes(c *gin.Context) {
	products, err := env.repo.GetProducts("fishes")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "base.html", gin.H{
			"error": http.StatusInternalServerError,
		})
		return
	}

	c.HTML(http.StatusOK, "base.html", gin.H{
		"title":    "Fishes",
		"Products": products,
	})
}

// GetDrinks retrieves all drink products from the repository and renders the view.
func (env *env) GetDrinks(c *gin.Context) {
	products, err := env.repo.GetProducts("drinks")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "base.html", gin.H{
			"error": http.StatusInternalServerError,
		})
		return
	}

	c.HTML(http.StatusOK, "base.html", gin.H{
		"title":    "Drinks",
		"Products": products,
	})
}

// SignupGET renders the signup page for user registration.
func (env *env) SignupGET(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", nil)
}

// SignupPOST handles the form submission for user registration.
func (env *env) SignupPOST(c *gin.Context) {
	var signupForm SignupForm

	if err := c.ShouldBind(&signupForm); err != nil {
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{
			"error": http.StatusText(http.StatusBadRequest),
		})
		return
	}

	// Validate the full name.
	err := ValidateFullName(signupForm.FullName)
	if err != nil {
		if errors.Is(err, ErrInvalidFullName) {
			c.HTML(http.StatusBadRequest, "signup.html", gin.H{
				"error_full_name": "Please enter a valid full name",
			})
		}
		return
	}

	// Sanitize and validate phone number.
	sanitizedPhoneNumber, err := sanitizePhoneNumber(signupForm.PhoneNumber)
	if err != nil {
		if errors.Is(err, ErrInvalidPhoneNumber) {
			c.HTML(http.StatusBadRequest, "signup.html", gin.H{
				"error_phone": "Invalid phone number",
			})
		} else {
			c.HTML(http.StatusBadRequest, "signup.html", gin.H{
				"error": http.StatusText(http.StatusBadRequest),
			})
		}
		return
	}

	// Check password strength.
	err = CheckPasswordStrength(signupForm.Password)
	if err != nil {
		if errors.Is(err, ErrWeakPassword) {
			c.HTML(http.StatusUnprocessableEntity, "signup.html", gin.H{
				"error_password": "Password must be 8+ characters with uppercase, lowercase, digit, and special character.",
			})
		}
		return
	}

	// Insert user data into the repository.
	err = env.repo.InsertUsers(signupForm.FullName, sanitizedPhoneNumber, signupForm.Password)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicatePhoneNumber) {
			c.HTML(http.StatusUnprocessableEntity, "signup.html", gin.H{
				"error_user": "User already exist",
			})
		} else {
			c.HTML(http.StatusInternalServerError, "signup.html", gin.H{
				"error": http.StatusText(http.StatusInternalServerError),
			})
		}
		return
	}

	c.Redirect(http.StatusFound, "/login")
}

// LoginGET renders the login page.
func (env *env) LoginGET(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

// LoginPOST handles the login form submission, including validation and authentication.
func (env *env) LoginPOST(c *gin.Context) {
	var loginForm LoginForm

	if err := c.ShouldBind(&loginForm); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": http.StatusText(http.StatusBadRequest),
		})
		return
	}

	// Sanitize and validate phone number.
	sanitizedPhoneNumber, err := sanitizePhoneNumber(loginForm.PhoneNumber)
	if err != nil {
		if errors.Is(err, ErrInvalidPhoneNumber) {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"error_phone": "Invalid phone number",
			})
		} else {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"error": http.StatusText(http.StatusBadRequest),
			})
		}
		return
	}

	// Authenticate the user.
	userID, err := env.repo.Authenticate(sanitizedPhoneNumber, loginForm.Password)
	if err != nil {
		if errors.Is(err, repository.ErrInvalidCredentials) {
			c.HTML(http.StatusUnprocessableEntity, "login.html", gin.H{
				"error_invalid_creds": "Invalid user credentials",
			})

		} else {
			c.HTML(http.StatusInternalServerError, "login.html", gin.H{
				"error": http.StatusText(http.StatusInternalServerError),
			})
		}
		return
	}

	// Create and sign the JWT token for the user.
	tokenString, err := env.createAndSignJWT(userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{
			"error": http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	// Set the JWT token in a cookie and redirect to the user account page.
	SetCookie(c, tokenString)
	c.Redirect(http.StatusSeeOther, "/user/orders")
}

// UserAccount renders the user account page.
func (env *env) UserAccount(c *gin.Context) {
	userID, _ := c.Get("userID")
	uID := int(userID.(float64))

	userInfo, err := env.repo.GetUserInfoByID(uID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "account.html", gin.H{
			"error": http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	c.HTML(http.StatusOK, "account.html", gin.H{
		"title":        "Account",
		"full_name":    userInfo.FullName,
		"phone_number": userInfo.PhoneNumber,
	})
}

// UpdateUserName handles the form submission for updating the user's name.
func (env *env) UpdateUserName(c *gin.Context) {
	var userUpdateForm UserUpdateForm

	if err := c.ShouldBind(&userUpdateForm); err != nil {
		c.HTML(http.StatusBadRequest, "account.html", gin.H{
			"error": http.StatusText(http.StatusBadRequest),
		})
		return
	}

	userID, _ := c.Get("userID")
	uID := int(userID.(float64))

	// Update the user's full name in the repository.
	if err := env.repo.UpdateUserNameByID(uID, userUpdateForm.FullName); err != nil {
		c.HTML(http.StatusInternalServerError, "account.html", gin.H{
			"error": http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/user/account")
}

// LogoutPOST logs out the user by deleting the authentication cookie.
func (env *env) LogoutPOST(c *gin.Context) {
	c.SetCookie("Auth", "deleted", 0, "/", "", false, true)

	c.Redirect(http.StatusSeeOther, "/")
}

// MakeAnOrder processes the user's order, calculates the total price, and saves the order in the database.
func (env *env) MakeAnOrder(c *gin.Context) {
	var order BasketOrder

	if err := c.ShouldBind(&order); err != nil {
		c.HTML(http.StatusBadRequest, "base.html", gin.H{
			"error": http.StatusText(http.StatusBadRequest),
		})
		return
	}

	// Calculate the total price of the order.
	var totalPrice float64
	for _, bItem := range order.OrderProducts {
		p, err := env.repo.GetProductByID(bItem.ID)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "base.html", gin.H{
				"error": http.StatusText(http.StatusInternalServerError),
			})
			return
		}
		subtotal := float64(p.PPrice) * float64(bItem.Quantity)
		totalPrice += subtotal
	}

	userID, _ := c.Get("userID")
	uID := int(userID.(float64))

	// Create a new order in the database.
	orderID, err := env.repo.CreateNewOrder(uID, totalPrice, order.DeliveryAddress)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "base.html", gin.H{
			"error": http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	// Add the order products to the order.
	for _, bItem := range order.OrderProducts {
		err := env.repo.CreateNewOrderProducts(orderID, bItem.ID, bItem.Quantity)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "base.html", gin.H{
				"error": http.StatusText(http.StatusInternalServerError),
			})
			return
		}
	}

	c.Redirect(http.StatusSeeOther, "/user/orders")
}

// Orders displays the user's previous orders.
func (env *env) Orders(c *gin.Context) {
	userID, _ := c.Get("userID")
	uID := int(userID.(float64))

	// Retrieve the user's orders from the repository.
	orders, err := env.repo.GetOrdersWithItems(uID)
	if err != nil {
		c.HTML(http.StatusOK, "orders.html", gin.H{
			"title": "Orders",
			"error": http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	// Sort the orders in descending order by OrderID
	sort.SliceStable(orders, func(i, j int) bool {
		return orders[i].OrderID > orders[j].OrderID
	})

	c.HTML(http.StatusOK, "orders.html", gin.H{
		"title":  "Orders",
		"orders": orders,
	})
}
