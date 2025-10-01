# DevOps Multi-API Project

A complete full-stack application demonstrating API development with Laravel, Go, and React frontend.

## 🏗️ Project Structure

```
devops/
├── laravel/          # Laravel API (PHP)
├── go/              # Go API 
├── frontend/        # React.js Frontend
├── compare_apis.sh  # API comparison script
└── README.md        # This file
```

## 🚀 How to Run Applications

### 🐘 Laravel API (Port 8001)
```bash
cd laravel
composer install
cp .env.example .env
php artisan key:generate
php artisan migrate --seed
php artisan serve --port=8001
```

### 🐹 Go API (Port 8080)
```bash
cd go
cp .env.example .env
go mod tidy
go run main.go
```

### ⚛️ React Frontend (Port 3000)
```bash
cd frontend
npm install
npm start
```

## 🌐 API Endpoints

Both Laravel and Go APIs provide identical endpoints:

- **GET** `/api/products` - Get all products
- **GET** `/api/products/{id}` - Get product by ID

## 🎯 Frontend Features

The React frontend provides:
- 🐘 **Laravel API Button** - Hits Laravel API at `http://127.0.0.1:8001`
- 🐹 **Go API Button** - Hits Go API at `http://localhost:8080`  
- 📊 **Real-time Data Display** - Shows API responses in formatted cards
- ❌ **Error Handling** - User-friendly error messages
- 📱 **Responsive Design** - Works on desktop and mobile

## 🔗 URLs

- **Laravel API:** http://127.0.0.1:8001/api/products
- **Go API:** http://localhost:8080/api/products  
- **React Frontend:** http://localhost:3000

## 🛠️ Tech Stack

- **Backend:** Laravel (PHP) + Go
- **Frontend:** React.js
- **Database:** MySQL (with SQLite fallback)
- **Styling:** CSS3 with gradients and animations

## ✅ Quick Test

Run all applications and test:

```bash
# Terminal 1: Laravel API
cd laravel && php artisan serve --port=8001

# Terminal 2: Go API  
cd go && go run main.go

# Terminal 3: React Frontend
cd frontend && npm start

# Terminal 4: Test APIs
./compare_apis.sh
```

Access frontend at http://localhost:3000 and click both API buttons to see the results!