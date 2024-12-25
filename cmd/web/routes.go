package main

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the application routes
func (env *env) SetupRoutes(router *gin.Engine) {
	// Public routes
	router.GET("/", env.GetBurgers)        // View burgers menu
	router.GET("/burgers", env.GetBurgers) // View burgers menu
	router.GET("/fishes", env.GetFishes)   // View fish menu
	router.GET("/drinks", env.GetDrinks)   // View drinks menu

	// Authentication routes
	router.GET("/signup", env.SignupGET)   // Display signup form
	router.POST("/signup", env.SignupPOST) // Handle signup form submission
	router.GET("/login", env.LoginGET)     // Display login form
	router.POST("/login", env.LoginPOST)   // Handle login form submission

	// Protected routes (authentication required)
	router.GET("/user/account", env.AuthMiddleware, env.UserAccount)     // View user account info
	router.PATCH("/user/update", env.AuthMiddleware, env.UpdateUserName) // Update user account info
	router.POST("/user/logout", env.AuthMiddleware, env.LogoutPOST)      // Logout the user

	// Order routes
	router.POST("/order/purchase", env.AuthMiddleware, env.MakeAnOrder) // Make a purchase
	router.GET("/user/orders", env.AuthMiddleware, env.Orders)          // View user orders

}
