package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"
)

// Product represents a product in the database
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Category    string  `json:"category"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// APIResponse represents the standard API response format
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Count   int         `json:"count,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Database connection
var db *sql.DB

// Database configuration functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDBConfig() (driver, user, password, name, host, port string) {
	driver = getEnv("DB_DRIVER", "mysql")
	user = getEnv("DB_USER", "root")
	password = getEnv("DB_PASSWORD", "")
	name = getEnv("DB_NAME", "laravel")
	host = getEnv("DB_HOST", "127.0.0.1")
	port = getEnv("DB_PORT", "3306")
	return
}

// Initialize database connection
func initDB() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default environment variables")
	}

	var err error
	driver, user, password, name, host, port := getDBConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, name)
	
	log.Printf("Connecting to database: %s@%s:%s/%s", user, host, port, name)
	
	db, err = sql.Open(driver, dsn)
	if err != nil {
		log.Printf("Error opening database: %v", err)
		log.Println("Continuing without database connection...")
		return
	}

	err = db.Ping()
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		log.Println("Continuing without database connection...")
		db = nil
		return
	}

	log.Println("Successfully connected to MySQL database")
}

// Create sample data if database is connected
func createSampleData() {
	if db == nil {
		return
	}

	// Create products table if it doesn't exist
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS products (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		price DECIMAL(10,2) NOT NULL,
		quantity INT DEFAULT 0,
		category VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}

	// Check if products already exist
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if err != nil {
		log.Printf("Error checking product count: %v", err)
		return
	}

	if count > 0 {
		log.Println("Products already exist in database")
		return
	}

	// Insert sample products
	sampleProducts := []Product{
		{Name: "Laptop Pro 15", Description: "High-performance laptop with 16GB RAM and 512GB SSD", Price: 1299.99, Quantity: 25, Category: "Electronics"},
		{Name: "Wireless Headphones", Description: "Noise-cancelling wireless headphones with 30h battery life", Price: 199.99, Quantity: 50, Category: "Electronics"},
		{Name: "Coffee Maker", Description: "Programmable coffee maker with 12-cup capacity", Price: 89.99, Quantity: 15, Category: "Home & Kitchen"},
		{Name: "Running Shoes", Description: "Lightweight running shoes with excellent cushioning", Price: 129.99, Quantity: 30, Category: "Sports & Outdoors"},
		{Name: "Smartphone", Description: "Latest smartphone with 128GB storage and triple camera", Price: 699.99, Quantity: 40, Category: "Electronics"},
	}

	for _, product := range sampleProducts {
		insertQuery := `
		INSERT INTO products (name, description, price, quantity, category) 
		VALUES (?, ?, ?, ?, ?)`
		
		_, err := db.Exec(insertQuery, product.Name, product.Description, product.Price, product.Quantity, product.Category)
		if err != nil {
			log.Printf("Error inserting product %s: %v", product.Name, err)
		}
	}

	log.Println("Sample products inserted successfully")
}

// getAllProducts retrieves all products from database or returns mock data
func getAllProducts() ([]Product, error) {
	if db == nil {
		// Return mock data if no database connection
		return getMockProducts(), nil
	}

	query := `SELECT id, name, description, price, quantity, category, created_at, updated_at FROM products`
	rows, err := db.Query(query)
	if err != nil {
		return getMockProducts(), err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		var createdAt, updatedAt time.Time
		
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, 
			&product.Quantity, &product.Category, &createdAt, &updatedAt)
		if err != nil {
			return getMockProducts(), err
		}
		
		product.CreatedAt = createdAt.Format("2006-01-02T15:04:05.000000Z")
		product.UpdatedAt = updatedAt.Format("2006-01-02T15:04:05.000000Z")
		products = append(products, product)
	}

	return products, nil
}

// getProductByID retrieves a specific product by ID
func getProductByID(id int) (*Product, error) {
	if db == nil {
		// Return mock data if no database connection
		mockProducts := getMockProducts()
		for _, product := range mockProducts {
			if product.ID == id {
				return &product, nil
			}
		}
		return nil, fmt.Errorf("product not found")
	}

	query := `SELECT id, name, description, price, quantity, category, created_at, updated_at FROM products WHERE id = ?`
	row := db.QueryRow(query, id)

	var product Product
	var createdAt, updatedAt time.Time
	
	err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price, 
		&product.Quantity, &product.Category, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, err
	}
	
	product.CreatedAt = createdAt.Format("2006-01-02T15:04:05.000000Z")
	product.UpdatedAt = updatedAt.Format("2006-01-02T15:04:05.000000Z")

	return &product, nil
}

// getMockProducts returns sample data when database is not available
func getMockProducts() []Product {
	return []Product{
		{ID: 1, Name: "Laptop Pro 15", Description: "High-performance laptop with 16GB RAM and 512GB SSD", Price: 1299.99, Quantity: 25, Category: "Electronics", CreatedAt: "2025-07-25T08:01:39.000000Z", UpdatedAt: "2025-07-25T08:01:39.000000Z"},
		{ID: 2, Name: "Wireless Headphones", Description: "Noise-cancelling wireless headphones with 30h battery life", Price: 199.99, Quantity: 50, Category: "Electronics", CreatedAt: "2025-07-25T08:01:39.000000Z", UpdatedAt: "2025-07-25T08:01:39.000000Z"},
		{ID: 3, Name: "Coffee Maker", Description: "Programmable coffee maker with 12-cup capacity", Price: 89.99, Quantity: 15, Category: "Home & Kitchen", CreatedAt: "2025-07-25T08:01:39.000000Z", UpdatedAt: "2025-07-25T08:01:39.000000Z"},
		{ID: 4, Name: "Running Shoes", Description: "Lightweight running shoes with excellent cushioning", Price: 129.99, Quantity: 30, Category: "Sports & Outdoors", CreatedAt: "2025-07-25T08:01:39.000000Z", UpdatedAt: "2025-07-25T08:01:39.000000Z"},
		{ID: 5, Name: "Smartphone", Description: "Latest smartphone with 128GB storage and triple camera", Price: 699.99, Quantity: 40, Category: "Electronics", CreatedAt: "2025-07-25T08:01:39.000000Z", UpdatedAt: "2025-07-25T08:01:39.000000Z"},
	}
}

// Handler for GET /api/products
func productsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response := APIResponse{
			Success: false,
			Message: "Method not allowed",
			Error:   "Only GET method is allowed",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	products, err := getAllProducts()
	if err != nil {
		response := APIResponse{
			Success: false,
			Message: "Error retrieving products",
			Error:   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := APIResponse{
		Success: true,
		Message: "Products retrieved successfully",
		Data:    products,
		Count:   len(products),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Handler for GET /api/products/{id}
func productHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response := APIResponse{
			Success: false,
			Message: "Method not allowed",
			Error:   "Only GET method is allowed",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(path)
	if err != nil {
		response := APIResponse{
			Success: false,
			Message: "Invalid product ID",
			Error:   "Product ID must be a number",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	product, err := getProductByID(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "product not found" {
			statusCode = http.StatusNotFound
		}

		response := APIResponse{
			Success: false,
			Message: "Product not found",
			Error:   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := APIResponse{
		Success: true,
		Message: "Product retrieved successfully",
		Data:    product,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Handler for root path
func homeHandler(w http.ResponseWriter, r *http.Request) {
	response := APIResponse{
		Success: true,
		Message: "Go Products API is running",
		Data: map[string]string{
			"endpoints": "GET /api/products, GET /api/products/{id}",
			"database":  func() string { if db != nil { return "MySQL connected" } else { return "Using mock data" } }(),
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Initialize database connection
	initDB()
	
	// Create sample data if database is connected
	createSampleData()

	// Set up routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api/products", productsHandler)
	http.HandleFunc("/api/products/", productHandler) // Handle /api/products/{id}

	log.Println("Go Products API Server is running on http://localhost:8080")
	log.Println("Endpoints:")
	log.Println("  GET /api/products     - Get all products")
	log.Println("  GET /api/products/{id} - Get product by ID")
	
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
