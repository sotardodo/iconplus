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
	CreatedAt   *string `json:"created_at"`
	UpdatedAt   *string `json:"updated_at"`
}

// APIResponse represents the standard API response format
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Count   int         `json:"count,omitempty"`
	Error   string      `json:"error,omitempty"`
}

var db *sql.DB

// Middleware CORS
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

// Load env or fallback
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Get database configuration
func getDBConfig() (driver, user, password, name, host, port string) {
	driver = getEnv("DB_DRIVER", "mysql")
	user = getEnv("DB_USER", "sotar")
	password = getEnv("DB_PASSWORD", "sotar123")
	name = getEnv("DB_NAME", "laravel")
	host = getEnv("DB_HOST", "127.0.0.1")
	port = getEnv("DB_PORT", "3306")
	return
}

// Initialize DB
func initDB() {
	godotenv.Load()

	var err error
	driver, user, password, name, host, port := getDBConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, name)

	log.Printf("Connecting to DB: %s@%s:%s/%s", user, host, port, name)

	db, err = sql.Open(driver, dsn)
	if err != nil || db.Ping() != nil {
		log.Println("⚠ DB not connected — using mock data")
		db = nil
		return
	}

	log.Println("✅ Connected to MySQL")
}

// Sample products creation
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
		created_at TIMESTAMP NULL DEFAULT NULL,
		updated_at TIMESTAMP NULL DEFAULT NULL
	)`

	_, _ = db.Exec(createTableQuery)
}

// Get all products (✅ FIXED NULL TIMESTAMP BUG)
func getAllProducts() ([]Product, error) {
	if db == nil {
		return getMockProducts(), nil
	}

	rows, err := db.Query(`SELECT id, name, description, price, quantity, category, created_at, updated_at FROM products`)
	if err != nil {
		return getMockProducts(), err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		var createdAt, updatedAt sql.NullTime

		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price,
			&product.Quantity, &product.Category, &createdAt, &updatedAt)
		if err != nil {
			continue
		}

		if createdAt.Valid {
			t := createdAt.Time.Format("2006-01-02T15:04:05")
			product.CreatedAt = &t
		}

		if updatedAt.Valid {
			t := updatedAt.Time.Format("2006-01-02T15:04:05")
			product.UpdatedAt = &t
		}

		products = append(products, product)
	}

	return products, nil
}

func getMockProducts() []Product {
	return []Product{
		{ID: 1, Name: "Laptop Pro 15", Price: 1299.99, Quantity: 25, Category: "Electronics"},
	}
}

// Routes
func productsHandler(w http.ResponseWriter, r *http.Request) {
	products, err := getAllProducts()
	if err != nil {
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Error retrieving products", Error: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(APIResponse{Success: true, Data: products, Count: len(products)})
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, _ := strconv.Atoi(idStr)

	products, _ := getAllProducts()
	for _, p := range products {
		if p.ID == id {
			json.NewEncoder(w).Encode(APIResponse{Success: true, Data: p})
			return
		}
	}

	json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Product not found"})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Message: "Go Products API is running",
	})
}

// Main
func main() {
	initDB()
	createSampleData()

	http.HandleFunc("/", withCORS(homeHandler))
	http.HandleFunc("/api/products", withCORS(productsHandler))
	http.HandleFunc("/api/products/", withCORS(productHandler))

	log.Println("🚀 Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
