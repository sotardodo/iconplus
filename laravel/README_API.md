# Laravel Products API

This Laravel application provides a REST API for managing products with JSON responses.

## API Endpoints

### GET /api/products
Retrieves all products from the database.

**Response:**
```json
{
    "success": true,
    "message": "Products retrieved successfully",
    "data": [
        {
            "id": 1,
            "name": "Laptop Pro 15",
            "description": "High-performance laptop with 16GB RAM and 512GB SSD",
            "price": "1299.99",
            "quantity": 25,
            "category": "Electronics",
            "created_at": "2025-07-25T08:01:39.000000Z",
            "updated_at": "2025-07-25T08:01:39.000000Z"
        }
    ],
    "count": 5
}
```

### GET /api/products/{id}
Retrieves a specific product by ID.

**Response:**
```json
{
    "success": true,
    "message": "Product retrieved successfully",
    "data": {
        "id": 1,
        "name": "Laptop Pro 15",
        "description": "High-performance laptop with 16GB RAM and 512GB SSD",
        "price": "1299.99",
        "quantity": 25,
        "category": "Electronics",
        "created_at": "2025-07-25T08:01:39.000000Z",
        "updated_at": "2025-07-25T08:01:39.000000Z"
    }
}
```

## Database Setup

### Current Setup (SQLite)
The application is currently configured to use SQLite for easy testing. The database file is located at `database/database.sqlite`.

### Switching to MySQL

To use MySQL instead of SQLite:

1. **Install and start MySQL server**
2. **Create a database:**
   ```sql
   CREATE DATABASE laravel;
   ```

3. **Update the `.env` file:**
   ```env
   DB_CONNECTION=mysql
   DB_HOST=127.0.0.1
   DB_PORT=3306
   DB_DATABASE=laravel
   DB_USERNAME=your_mysql_username
   DB_PASSWORD=your_mysql_password
   ```

4. **Run migrations and seeders:**
   ```bash
   php artisan migrate:fresh --seed
   ```

## Running the Application

1. **Start the development server:**
   ```bash
   php artisan serve
   ```

2. **Test the API endpoints:**
   ```bash
   # Get all products
   curl -H "Accept: application/json" http://127.0.0.1:8000/api/products
   
   # Get specific product
   curl -H "Accept: application/json" http://127.0.0.1:8000/api/products/1
   ```

## Files Created/Modified

- **Model:** `app/Models/Product.php` - Product model with fillable fields
- **Controller:** `app/Http/Controllers/ProductController.php` - API controller with JSON responses
- **Migration:** `database/migrations/2025_07_25_075615_create_products_table.php` - Products table structure
- **Seeder:** `database/seeders/ProductSeeder.php` - Sample product data
- **Routes:** `routes/api.php` - API route definitions

## Database Schema

The `products` table includes:
- `id` (primary key)
- `name` (string)
- `description` (text, nullable)
- `price` (decimal 10,2)
- `quantity` (integer, default 0)
- `category` (string, nullable)
- `created_at` and `updated_at` (timestamps)

## Error Handling

The API includes proper error handling:
- Returns JSON responses for all endpoints
- Includes success/error status
- Provides meaningful error messages
- Uses appropriate HTTP status codes (200, 404, 500)
