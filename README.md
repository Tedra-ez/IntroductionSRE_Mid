# Clothes Store

A premium, modern Online store platform built with Go and MongoDB. Features a luxury design aesthetic with a fully functional shopping experience and a comprehensive administrative dashboard.

## Team
- Yskak Zhanibek
- Nauanov Alikhan
- Zhumagali Beibarys

## Features

### Customer Experience
- **Modern Shop**: Advanced filtering (category, gender, color, size) and sorting.
- **Product Details**: High-quality imagery, size selection, and stock status.
- **Cart & Wishlist**: Persistent client-side shopping cart and wishlist management.
- **Checkout**: Seamless checkout flow with address management and order confirmation.
- **User Accounts**: Registration, login, and order history tracking.

### Admin Dashboard
- **Analytics**: Key performance indicators (Total Sales, Orders, Users).
- **Product Management**: Complete CRUD with local image uploads and advanced validation.
- **Order Management**: Track and update order statuses.
- **User Management**: Overview of registered users.

## Tech Stack

- **Backend**: Go (Gin Web Framework)
- **Database**: MongoDB
- **Authentication**: JWT (JSON Web Tokens) with Secure Cookies
- **Frontend**: Semantic HTML5, Vanilla CSS (Modern CSS variables), JavaScript (ES6+)
- **Icons**: Lucide Icons

## Database Performance

- **Multi-stage aggregation**: Analytics uses MongoDB pipelines (`$facet`, `$group`, `$lookup`, `$sort`) to compute totals, revenue trends, and top products without loading every order into memory.
- **Compound indexes**: `orders` uses `{ userId: 1, createdAt: -1 }` for user history and recent sorting; `order_items` uses `{ orderId: 1, productId: 1 }` to accelerate joins and product sales grouping.
- **Reduced transfer**: Aggregations return compact summaries and only a small window of recent orders.

## Project Structure

```bash
├── cmd/
│   └── server/          # Entry point (main.go)
├── internal/
│   ├── api/             # Routing and Middleware
│   ├── config/          # Environment configuration
│   ├── db/              # Database connection
│   ├── handlers/        # HTTP Handlers
│   ├── models/          # Data structures
│   ├── repository/      # Database operations
│   └── services/        # Business logic
├── static/
│   ├── assets/          # Images, Banners, UI elements
│   ├── css/             # Stylesheets
│   └── js/              # Client-side logic
└── templates/           # HTML fragments
```

## Setup & Installation

### Prerequisites
- Go 1.25+
- MongoDB instance

### 1. Environment Configuration
Create a `.env` file in the root directory:
```env
PORT=8000
MONGODB_URI=Nelzya
JWT_SECRET=No
ADMIN_EMAIL=No
```

### 2. Run the Application
```bash
go run cmd/server/main.go
```
The server will start at http://localhost:8000.

## API Documentation

### Authentication
- **POST** `/auth/register`
  - Request (JSON):
    ```json
    { "fullName": "John Doe", "email": "john@example.com", "password": "secret123" }
    ```
  - Response `201`:
    ```json
    { "message": "user registered" }
    ```
- **POST** `/auth/login`
  - Request (JSON):
    ```json
    { "email": "john@example.com", "password": "secret123" }
    ```
  - Response `200`:
    ```json
    { "token": "jwt-token" }
    ```
- **GET** `/auth/logout` (browser redirect)

**Auth for API**: send `Authorization: Bearer <token>` or cookie `auth_token`.

### Products
- **GET** `/api/product`
  - Response `200`:
    ```json
    [{ "id": "p1", "name": "Sneakers", "price": 120, "sizes": ["41","42"], "colors": ["black"] }]
    ```
- **GET** `/api/product/:id`
  - Response `200`: product object
- **POST** `/api/product` (admin, multipart/form-data)
  - Example:
    ```bash
    curl -X POST http://localhost:8000/api/product \
      -H "Authorization: Bearer <token>" \
      -F "name=Sneakers" -F "price=120" -F "category=shoes" \
      -F "gender=unisex" -F "sizes=41,42" -F "colors=black" \
      -F "stock=41:5,42:3" -F "image=@./sneakers.jpg"
    ```
  - Response `201`: product object
- **PUT** `/api/product/:id` (admin, JSON body = product)
- **DELETE** `/api/product/:id` (admin) → `204`

### Orders
- **GET** `/orders?user_id={userId}`
  - Response `200`: array of orders
- **POST** `/orders`
  - Request (JSON):
    ```json
    {
      "user_id": "u1",
      "payment_method": "card",
      "delivery_method": "courier",
      "delivery_address": "Almaty, Abay 10",
      "comment": "leave at door",
      "items": [
        {
          "product_id": "p1",
          "product_name": "Sneakers",
          "selected_size": "42",
          "selected_color": "black",
          "quantity": 1,
          "unit_price": 120
        }
      ]
    }
    ```
  - Response `201`: order object
- **GET** `/orders/:id`
  - Response `200`: order object
- **PATCH** `/orders/:id/status`
  - Request (JSON):
    ```json
    { "status": "completed" }
    ```
  - Response `200`:
    ```json
    { "id": "orderId", "status": "completed" }
    ```

### Analytics (admin)
- **GET** `/api/analytics/stats` → dashboard stats
- **GET** `/api/analytics/top-products` → top product sales
- **GET** `/api/analytics/revenue?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD`
  - Response `200`:
    ```json
    [{ "date": "2026-02-01", "revenue": 520, "orders": 4 }]
    ```
- **GET** `/api/analytics/orders-status` → `{ "pending": 2, "completed": 5 }`

## Code Quality
- **Clean Architecture**: Separation of concerns between layers.
- **Optimized Assets**: Localized assets for faster loading and reliability.
- **Sanitized**: Codebase is free of redundant comments and junk files.
