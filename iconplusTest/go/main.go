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

// 🔹 CORS ADDED: Middleware sederhana
func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	}
}

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

	sampleProducts := []Product{
		{Name: "Laptop Pro 15", Description: "High-performance laptop", Price: 1299.99, Quantity: 25, Category: "Electronics"},
		{Name: "Wireless Headphones", Description: "Noise-cancelling wireless headphones", Price: 199.99, Quantity: 50, Category: "Electronics"},
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

func getMockProducts() []Product {
	return []Product{
		{ID: 1, Name: "Laptop Pro 15", Price: 1299.99, Quantity: 25, Category: "Electronics"},
	}
}

// Handlers...
func productsHandler(w http.ResponseWriter, r *http.Request) {
	// ... kode lama tetap
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	// ... kode lama tetap
}

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
	initDB()
	createSampleData()

	// 🔹 CORS ADDED: Bungkus semua handler dengan middleware
	http.HandleFunc("/", withCORS(homeHandler))
	http.HandleFunc("/api/products", withCORS(productsHandler))
	http.HandleFunc("/api/products/", withCORS(productHandler))

	log.Println("Go Products API Server is running on http://localhost:8080")
	log.Println("Endpoints:")
	log.Println("  GET /api/products     - Get all products")
	log.Println("  GET /api/products/{id} - Get product by ID")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
