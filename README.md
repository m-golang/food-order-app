# Food Order Web App

A simple and secure web application for ordering food. Users can sign up, log in, view the menu, place orders, and track their order history. This app features secure authentication with JWT, input validation, and a structured database.

## Features

### User Authentication:

- Signup and login with phone number and password.
- JWT-based authentication with secure cookies.
- Password strength validation.

### Order Management:

- Browse product menus: Burgers, Fishes, Drinks.
- Place orders with selected products.
- View past orders with itemized details.

### User Profile:

- Update full name in user profile.
- View and manage user orders.

### Security:

- JWT tokens stored as HTTP-only cookies.
- Phone number and password validation.
- Password encryption using bcrypt.

### Backend:

- Built with Go using the Gin web framework.
- Data stored in a MySQL database.

### Frontend:

- Dynamic HTML templates served with Gin.

## Technologies Used

- **Go (Golang):** Backend server and API development.
- **Gin:** A web framework for building the application.
- **MySQL:** Relational database to store user data, orders, and products.
- **JWT:** JSON Web Tokens for secure authentication.
- **bcrypt:** For password hashing and security.
- **HTML/CSS:** User interface rendered using templates.

## Installation

### Prerequisites

- To run this app, you'll need the following:

- **Go:** Version 1.18 or higher.
- **MySQL:** A MySQL server running locally.
- **Git:** To clone the repository.

## Steps to Run the Application Locally

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/m-golang/food-order-app.git 
    cd food-order-app

2.  **Set up the MySQL Database:**

    - Create a new database in MySQL for the app. For example:

        ```bash
        CREATE DATABASE burgerfish;
        
        CREATE TABLE `menu` (
            `id` INT NOT NULL AUTO_INCREMENT,
            `name` VARCHAR(50) NOT NULL,
            `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
            `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            PRIMARY KEY (`id`)
        );
        
        CREATE TABLE `products` (
            `id` INT NOT NULL AUTO_INCREMENT,
            `product_name` VARCHAR(100) NOT NULL,
            `product_description` TEXT,
            `product_price` DECIMAL(10,2) NOT NULL,
            `product_image` VARCHAR(255) DEFAULT NULL,
            `menu_id` INT NOT NULL,
            `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
            `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            PRIMARY KEY (`id`),
            CONSTRAINT `fk_menu` FOREIGN KEY (`menu_id`) REFERENCES `menu` (`id`) ON DELETE CASCADE
        );
        
        CREATE TABLE `orders` (
            `id` INT NOT NULL AUTO_INCREMENT,
            `user_id` INT NOT NULL,
            `order_date` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
            `total_amount` DECIMAL(10,2) NOT NULL,
            `delivery_address` TEXT NOT NULL,
            `status` ENUM('pending','completed','shipped','cancelled') DEFAULT 'pending',
            `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
            `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            PRIMARY KEY (`id`),
            CONSTRAINT `fk_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
        );
        
        CREATE TABLE `order_products` (
            `order_id` INT NOT NULL,
            `product_id` INT NOT NULL,
            `quantity` INT NOT NULL,
            PRIMARY KEY (`order_id`, `product_id`),
            CONSTRAINT `fk_order` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON DELETE CASCADE,
            CONSTRAINT `fk_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE
        );
        
        CREATE TABLE `users` (
            `id` INT NOT NULL AUTO_INCREMENT,
            `full_name` VARCHAR(50) NOT NULL,
            `phone_number` VARCHAR(14) NOT NULL,
            `password_hash` CHAR(60) DEFAULT NULL,
            `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
            `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            PRIMARY KEY (`id`),
            UNIQUE KEY `phone_number` (`phone_number`)
        );
        
        -- Creating indexes for faster queries
        
        CREATE INDEX idx_users_phone_number ON users (phone_number);       -- Index on phone number for fast lookup
        CREATE INDEX idx_orders_user_id ON orders (user_id);               -- Index on user_id for fast lookup
        CREATE INDEX idx_orders_status ON orders (status);                 -- Index on order status
        CREATE INDEX idx_orders_order_date ON orders (order_date);         -- Index on order date
        CREATE INDEX idx_orders_total_amount ON orders (total_amount);     -- Index on total amount for filtering
        CREATE INDEX idx_products_product_name ON products (product_name); -- Index on product name for searching


    - Make sure to adjust the database connection string in the `main.go` file (`dsn` variable).

3.  **Install Dependencies:** Ensure Go modules are set up and the necessary dependencies are installed:

    ```bash
    go mod tidy

4. **Configure Environment Variables:**

- Set the `dsn` (Data Source Name) for MySQL in the `main.go` file or set it as an environment variable:
    ```bash
    dsn := flag.String("dsn", "dbusername:dbpassword@/burgerfish?parseTime=true", "MySQL data source name")

5. **Run the Application:**
    ```bash
    go run cmd/web/main.go

## Routes

- **GET /signup:** Signup page.
- **POST /signup:** Create a new user.
- **GET /login:** Login page.
- **POST /login:** Authenticate and log in a user.
- **GET /burgers:** View the burgers menu.
- **GET /fishes:** View the fishes menu.
- **GET /drinks:** View the drinks menu.
- **POST /order/purchase:** Place an order.
- **GET /user/orders:** View user’s past orders.
- **PATCH /user/update:** Update the user's full name.
- **POST /user/logout:** Logout the user.
