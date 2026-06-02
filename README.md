# Product API

REST API service for managing (Create and update) products built with Golang, PostgreSQL, Clean Architecture, and Dependency Injection.

## Features

* Create Product API
* Patch Product API (Partial Update)
* PostgreSQL Database
* Clean Architecture
* Dependency Injection
* Swagger Documentation
* Unit Test
* Repository Integration Test
* Component Test
* Docker Compose Support

---

## Project Structure

```text
product-api
├── cmd                               // Application entry point
│   └── server
├── config                            // Configuration management
├── docs                              // Swagger documentation
├── internal
│   ├── domain                        // Core business entities
│   ├── dto                           // Request and response models
│   ├── handler                       // HTTP request handlers
│   ├── mapper                        // Data mapping
│   ├── repository                    // Database operations
│   ├── router                        // Route registration
│   └── service                       // Business logic
├── migrations                        // Database migrations
├── pkg                               // Shared utilities
│   ├── database                      // Database connection utilities
│   └── response                      // Common API response helpers
├── tests
│   └── component                     // End-to-end component tests
├── .gitignore
├── config.docker.yaml.example      
├── config.yaml.example
├── docker-compose.yml                // Docker Compose configuration.
├── dockerfile                        // Docker image definition.
└── README.md
```

---

## Architecture

```text
HTTP Request
     │
     ▼
 Handler
     │
     ▼
 Service (Usecase)
     │
     ▼
 Repository
     │
     ▼
 PostgreSQL
```

The project follows Clean Architecture principles by separating responsibilities into different layers.

* Handler: HTTP layer
* Service: Business logic / Use cases
* Repository: Data access layer
* Domain: Core business entities

Dependency Injection is used to inject dependencies between layers.

---

## API Documentation

Swagger UI is available at:

```text
http://localhost:3000/api-docs/index.html
```

---

## Running with Docker

Start application and PostgreSQL:

```bash
docker compose up --build
```

API:

```text
http://localhost:3000
```

Swagger:

```text
http://localhost:3000/api-docs/index.html
```

---

## Running Locally

### 1. Start PostgreSQL

```bash
docker compose up -d postgres
```

### 2. Create configuration file

Copy example configuration:

```bash
cp config.yaml.example config.yaml
```

### 3. Run application

```bash
go run cmd/server/main.go
```

---

## Running Tests

Run all tests:

```bash
go test ./... -v
```

---

## Database

Database: PostgreSQL

Migration file:

```text
migrations/001_create_products_table.sql
```

---

## API Endpoints

### Create Product

```http
POST /product
```

Request:

```json
{
  "name": "Johnson's Baby",
  "description": "Top-To-Toe Hair&Body Bath",
  "sale_price": 1278,
  "price": 1680
}
```

Response:

```json
{
  "successful": true,
  "error_code": "",
  "data": {
    "name": "Johnson's Baby",
    "description": "Top-To-Toe Hair&Body Bath",
    "sale_price": 1278,
    "price": 1680
  }
}
```

---

### Patch Product

```http
PATCH /product/{id}
```

Supports partial update.

Example:

```json
{
  "name": "New Product Name"
}
```

or

```json
{
  "price": 1800,
  "sale_price": 1500
}
```

Response:

```json
{
  "successful": true,
  "error_code": ""
}
```

---

## Validation Rules

### Create Product

* name is required
* description is empty or blank, it will be stored as null
* price is required
* price must be greater than or equal 0
* sale_price must be lower than price
* sale_price must be greater than or equal 0

### Patch Product

* only provided fields are updated
* name cannot be null or blank
* description is empty or blank, it will be stored as null
* sale_price must be lower than price
* sale_price must be greater than or equal 0
* price cannot be null
* price must be greater than or equal 0

---

## Test Coverage

Implemented test categories:

* Service Unit Test
* Repository Integration Test
* Component Test (HTTP → Service → Repository)
* Service/domain unit test (function)

Run:

```bash
go test ./... -v
```

---
