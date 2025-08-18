# Go RESTful API Service with Express.js-like Structure

A Go service built with Gin framework following Express.js patterns for code organization with separate controllers and routes folders.

## Intended Usecase 
The orchestration service is intended to be used in the second approach to the bio auth project at WSO2. In this approach its planned to utilize the facial recognition capabilites offered by the facial recognition capable Hikvision camera which is connected to the HikCentral Professioanal application that is run in the server at WSO2.

The orchestartion service is suppose to recieve the recognized user details from the alarm triggered at HikCentral (Due to a facial recognition event taking place) and then using SCIM2 API endpoints of the WSO2 IDP (Asgardoe or IS) recieve the roles of the corresponding user. Then the roles will be checked and if the required access level is present the required door can be opened.

There also needs to be a way to communicate through which camera the facial data was submitted. Then the relevant door can be identified and once the role verifcation is done that door can be opened.

## Project Structure

```
orchestration-service-approach2/
├── main.go                     # Application entry point
├── controllers/                # Business logic handlers (like Express.js controllers)
│   ├── health_controller.go    # Health check controller
│   └── item_controller.go      # Item operations controller
├── routes/                     # Route definitions (like Express.js routes)
│   ├── routes.go              # Main routes setup
│   ├── health_routes.go       # Health routes
│   └── item_routes.go         # Item API routes
├── models/                     # Data models and business entities
│   └── item.go                # Item model and data store
├── go.mod                     # Go module definition
└── README.md                  # This file
```

## Architecture Pattern

This structure follows the **MVC (Model-View-Controller)** pattern similar to Express.js:

- **Models** (`models/`): Data structures and business logic
- **Controllers** (`controllers/`): Handle HTTP requests and responses
- **Routes** (`routes/`): Define URL patterns and link them to controllers

## Key Features

- **Separation of Concerns**: Clean separation between routing, business logic, and data models
- **Modular Design**: Easy to add new features by creating new controllers and routes
- **Scalable Structure**: Similar to Express.js, making it familiar for Node.js developers
- **RESTful API**: Following REST principles with proper HTTP methods and status codes

## API Endpoints

### Health Check
- `GET /health` - Check service health status

### Items API (Example endpoints)
- `GET /api/v1/items` - Get all items
- `GET /api/v1/items/:id` - Get item by ID
- `POST /api/v1/items` - Create new item
- `PUT /api/v1/items/:id` - Update item by ID
- `DELETE /api/v1/items/:id` - Delete item by ID

## Getting Started

### Prerequisites
- Go 1.25.0 or later

### Installation & Running

1. Install dependencies:
   ```bash
   go mod tidy
   ```

2. Run the service:
   ```bash
   go run main.go
   ```

3. Or build and run:
   ```bash
   go build -o service main.go
   ./service
   ```

The service will start on port 8080.

## Adding New Features

### 1. Create a Model (if needed)
```go
// models/user.go
package models

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

### 2. Create a Controller
```go
// controllers/user_controller.go
package controllers

import (
    "github.com/gin-gonic/gin"
    "your-module/models"
)

type UserController struct{}

func (uc *UserController) GetUsers(c *gin.Context) {
    // Implementation here
}

func (uc *UserController) CreateUser(c *gin.Context) {
    // Implementation here
}
```

### 3. Create Routes
```go
// routes/user_routes.go
package routes

import (
    "github.com/gin-gonic/gin"
    "your-module/controllers"
)

func SetupUserRoutes(router *gin.Engine) {
    userController := &controllers.UserController{}
    
    v1 := router.Group("/api/v1")
    {
        users := v1.Group("/users")
        {
            users.GET("", userController.GetUsers)
            users.POST("", userController.CreateUser)
        }
    }
}
```

### 4. Register Routes in Main Routes File
```go
// routes/routes.go
func SetupRoutes(router *gin.Engine) {
    SetupHealthRoutes(router)
    SetupItemRoutes(router)
    SetupUserRoutes(router)  // Add this line
}
```

## Example Usage

### Create an item
```bash
curl -X POST http://localhost:8080/api/v1/items \
  -H "Content-Type: application/json" \
  -d '{"name": "New Item", "description": "A new test item"}'
```

### Get all items
```bash
curl http://localhost:8080/api/v1/items
```

### Health check
```bash
curl http://localhost:8080/health
```

## Benefits of This Structure

1. **Familiar to Express.js developers**: Similar folder structure and patterns
2. **Maintainable**: Clear separation of concerns
3. **Scalable**: Easy to add new features without cluttering main.go
4. **Testable**: Controllers can be easily unit tested
5. **Team-friendly**: Multiple developers can work on different parts without conflicts

## Next Steps for Production

- Add database integration in models
- Implement middleware for authentication/authorization
- Add input validation and sanitization
- Add comprehensive logging
- Add configuration management
- Add unit and integration tests
- Add Docker containerization
- Add API documentation (Swagger/OpenAPI)