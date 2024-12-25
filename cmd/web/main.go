package main

import (
	"database/sql"
	"flag"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/m-golang/food-order-app/internals/repository"

	_ "github.com/go-sql-driver/mysql"
)

type env struct {
	repo      *repository.RepoModel
	secretKey string
}

func main() {
	// Parse command-line arguments to get the database DSN (Data Source Name)
	dsn := flag.String("dsn", "dbusername:dbpassword@/burgerfish?parseTime=true", "MySQL data source name")
	flag.Parse()

	// Open a connection to the MySQL database using the DSN
	db, err := openDB(*dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close() // Ensure the database connection is closed when done

	env := &env{
		repo:      &repository.RepoModel{DB: db}, // Repository model with DB
		secretKey: "Z2VudGVyYXRlX3NlY3JldF9rZXkxMjM0NTY3ODl9YTlhYmMzZGRiZjMwMTY5YWZkNzE3NzY2YWY=",
	}

	// Create a new Gin router
	router := gin.Default()

	// Load HTML templates and static files for the frontend
	router.LoadHTMLGlob("ui/html/*.html")
	router.Static("/static", "./ui/static")

	// Set up the application routes
	env.SetupRoutes(router)

	// Start the Gin web server
	router.Run()
}

// Opens a MySQL database connection using the provided DSN
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
