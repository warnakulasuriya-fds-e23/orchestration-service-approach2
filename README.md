# Go RESTful API Service with Express.js-like Structure

A Go service built with Gin framework following Express.js patterns for code organization with separate controllers and routes folders.

## Intended Usecase 
The orchestration service is intended to be used in the second approach to the bio auth project at WSO2. In this approach its planned to utilize the facial recognition capabilites offered by the facial recognition capable Hikvision camera which is connected to the HikCentral Professioanal application that is run in the server at WSO2.

The orchestartion service is suppose to recieve the recognized user details from the alarm triggered at HikCentral (Due to a facial recognition event taking place) and then using SCIM2 API endpoints of the WSO2 IDP (Asgardoe or IS) recieve the roles of the corresponding user. Then the roles will be checked and if the required access level is present the required door can be opened.

There also needs to be a way to communicate through which camera the facial data was submitted. Then the relevant door can be identified and once the role verifcation is done that door can be opened.

## Project Structure

```
orchestration-service-approach2/
├── main.go                               # Application entry point
├── .env                                  # Environment variables (not tracked in git)
├── internal/                             # Internal application code
│   ├── controllers/                      # Business logic handlers organized by version
│   │   ├── common_controllers/           # Common controllers across versions
│   │   │   └── health_controller.go      # Health check controller
│   │   └── v1_controllers/               # Version 1 specific controllers
│   │       └── authorization_controller.go # User authorization logic
│   ├── models/                           # Data models and business entities
│   │   └── incoming_data_from_hikcentral.go # HikCentral data structures
│   └── routes/                           # Route definitions organized hierarchically
│       ├── routes.go                     # Main routes setup
│       └── internal/                     # Internal route organization
│           ├── common_routes/            # Common routes across versions
│           │   └── health_routes.go      # Health check routes
│           └── v1_routes/                # Version 1 specific routes
│               └── authorization_routes.go # Authorization routes
├── controllers/                          # Legacy controllers (for backward compatibility)
├── routes/                              # Legacy routes (for backward compatibility)  
├── models/                              # Legacy models (for backward compatibility)
├── go.mod                               # Go module definition
├── go.sum                               # Go dependencies checksum
└── README.md                            # This file
```

## Architecture Pattern

This structure follows a **Clean Architecture** pattern with versioned API design:

- **Internal Package Structure**: All business logic is contained within the `internal/` package following Go best practices
- **Versioned Controllers**: Controllers are organized by API version (v1, v2, etc.) for backward compatibility
- **Common Components**: Shared functionality like health checks are in common packages
- **Models** (`internal/models/`): Data structures for external integrations (HikCentral, SCIM2)
- **Controllers** (`internal/controllers/`): Handle HTTP requests, business logic, and external API calls
- **Routes** (`internal/routes/`): Define URL patterns organized by version and functionality

## Key Features

- **Clean Architecture**: Internal package structure following Go best practices
- **Versioned API Design**: Support for multiple API versions (v1, v2, etc.)
- **HikCentral Integration**: Receives facial recognition data from HikCentral Professional
- **SCIM2 Integration**: Connects with WSO2 Identity Server for user role verification
- **Environment Configuration**: Uses .env files for secure configuration management
- **Modular Design**: Easy to add new features by creating new versioned controllers and routes
- **Scalable Structure**: Organized for enterprise-level facial recognition access control
- **RESTful API**: Following REST principles with proper HTTP methods and status codes



## API Endpoints

### Health Check
- `GET /health` - Check service health status

### Authorization API (v1)
- `POST /api/v1/authorize-for-door-access` - Authorize user for door access based on facial recognition data

#### Authorization Endpoint Details
This endpoint receives facial recognition data from HikCentral and performs the following:
1. Receives user identification data from HikCentral alarm
2. Queries WSO2 Identity Server using SCIM2 API to get user roles
3. Validates user permissions for door access
4. Returns authorization decision

**Request Body Example:**
```json
{
  "id": "user123",
  "name": "John Doe", 
  "description": "Facial recognition match"
}
```

## Getting Started

### Prerequisites
- Go 1.25.0 or later

# Environmental Variable setup
Following environmental variables should be available for the proper functioning
of the service.
``` 
IDP_ADDRESS=https://localhost:9443
IDP_USERNAME=admin
IDP_PASSWORD=admin
```

**Environment Variables Description:**
- `IDP_ADDRESS`: Base URL of the WSO2 Identity Server instance
- `IDP_USERNAME`: Username for SCIM2 API authentication 
- `IDP_PASSWORD`: Password for SCIM2 API authentication

Create a `.env` file in the project root with these variables for local development.

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

The service will start on port 5000.

## Adding New Features

### 1. Create a Model (if needed)
```go
// internal/models/new_model.go
package models

type NewModel struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Data  string `json:"data"`
}

func (n *NewModel) ProcessData() error {
    // Implementation here
    return nil
}
```

### 2. Create a Versioned Controller
```go
// internal/controllers/v1_controllers/new_controller.go
package v1_controllers

import (
    "github.com/gin-gonic/gin"
    "github.com/warnakulasuriya-fds-e23/orchestration-service-approach2/internal/models"
)

type NewController struct{}

func (nc *NewController) HandleNewRequest(c *gin.Context) {
    // Implementation here
}
```

### 3. Create Versioned Routes
```go
// internal/routes/internal/v1_routes/new_routes.go
package v1_routes

import (
    "github.com/gin-gonic/gin"
    "github.com/warnakulasuriya-fds-e23/orchestration-service-approach2/internal/controllers/v1_controllers"
)

func SetupNewRoutes(v1Group *gin.RouterGroup) {
    controller := v1_controllers.NewController{}
    v1Group.POST("/new-endpoint", controller.HandleNewRequest)
}
```

### 4. Register Routes in Main Routes File
```go
// internal/routes/routes.go
func SetupRoutes(router *gin.Engine) {
    common_routes.SetupHealthRoutes(router)
    
    v1 := router.Group("/api/v1")
    {
        v1_routes.SetupAuthorizationRoutes(v1)
        v1_routes.SetupNewRoutes(v1)  // Add this line
    }
}
```

## Example Usage

### Authorize user for door access
```bash
curl -X POST http://localhost:5000/api/v1/authorize-for-door-access \
  -H "Content-Type: application/json" \
  -d '{
    "id": "user123",
    "name": "John Doe",
    "description": "Facial recognition match from camera 1"
  }'
```

### Health check
```bash
curl http://localhost:5000/health
```

**Expected Response for Authorization:**
```json
{
  "id": "user123",
  "userName": "john.doe",
  "emails": [...],
  "groups": [...],
  "roles": [...]
}
```

## Benefits of This Structure

1. **Clean Architecture**: Follows Go best practices with internal package organization
2. **Versioned APIs**: Easy to maintain backward compatibility while adding new features
3. **Enterprise-Ready**: Structured for integration with enterprise identity systems
4. **Maintainable**: Clear separation of concerns across layers
5. **Scalable**: Easy to add new API versions and features without breaking existing functionality
6. **Secure**: Environment-based configuration for sensitive credentials
7. **Testable**: Controllers and models can be easily unit tested
8. **Team-friendly**: Multiple developers can work on different API versions simultaneously

## Integration Points

### HikCentral Professional
- Receives alarm notifications when facial recognition occurs
- Processes user identification data from camera systems
- Supports multiple camera configurations for different access points

### WSO2 Identity Server
- Uses SCIM2 protocol for user data retrieval
- Supports role-based access control
- Handles authentication via basic auth (configurable)

## Security Considerations

- TLS verification is currently disabled for development (should be enabled in production)
- Basic authentication used for SCIM2 API calls
- Environment variables used for credential management
- No API key authentication implemented yet (recommended for production)

## Next Steps for Production

- Enable TLS verification for SCIM2 API calls
- Implement API key authentication for HikCentral webhooks
- Add comprehensive input validation and sanitization
- Implement door control integration (hardware APIs)
- Add camera-to-door mapping configuration
- Add comprehensive logging and monitoring
- Add role-based authorization logic
- Add unit and integration tests
- Add Docker containerization
- Add API documentation (Swagger/OpenAPI)
- Implement rate limiting and request throttling
- Add database persistence for audit logs
- Configure production-ready CORS policies