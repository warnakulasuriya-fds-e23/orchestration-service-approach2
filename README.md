# Go RESTful API Service for Facial Recognition Access Control

A Go service built with the Gin framework, following Express.js-inspired patterns for code organization with separate controllers, models, and routes folders.

## Intended Usecase

This orchestration service is designed for the bio-auth project at WSO2. It integrates with Hikvision cameras via HikCentral Professional, receives facial recognition events, and queries WSO2 Identity Server (IDP) using SCIM2 API endpoints to retrieve user roles. Based on roles and access requirements, it determines if a user should be granted door access.

## Project Structure

```
orchestration-service-approach2/
├── main.go
├── .env
├── accessRequirementsForDevices.json      # Device access requirements mapping
├── internal/
│   ├── controllers/
│   │   ├── common_controllers/
│   │   │   ├── health_controller.go
│   │   │   └── logger_controller.go
│   │   └── v1_controllers/
│   │       ├── authorization_controller.go
│   │       └── requirements_controller.go
│   ├── models/
│   │   ├── access_requirements_object.go
│   │   ├── incoming_data_from_hikcentral.go
│   │   └── wso2_idp_scim_role_object.go
│   ├── routes/
│   │   ├── routes.go
│   │   └── internal/
│   │       ├── common_routes/
│   │       │   ├── health_routes.go
│   │       │   └── logger_routes.go
│   │       └── v1_routes/
│   │           ├── authorization_routes.go
│   │           └── requirements_routes.go
│   └── utils/
│       ├── readAccessRequirementsfile.go
│       ├── requirementsManager.go
│       └── role_based_authorization.go
├── builds/
│   ├── orchestration-service-approach2
│   ├── orchestration-service-approach2.zip
│   ├── orchestration.exe
│   ├── orchestration.zip
│   └── service
├── go.mod
├── go.sum
├── open-api-v1-spec.yaml                  # OpenAPI spec for v1 endpoints
└── README.md
```

## Architecture Pattern

- **Clean Architecture**: All business logic is in the `internal/` package.
- **Versioned Controllers & Routes**: Organized by API version for backward compatibility.
- **Common Components**: Shared logic (health, logging) in `common_controllers` and `common_routes`.
- **Models**: Data structures for HikCentral, SCIM2, and access requirements.
- **Utils**: Helper functions for requirements management and role-based authorization.

## Key Features

- **Facial Recognition Integration**: Receives and processes events from HikCentral.
- **Role-Based Access Control**: Uses SCIM2 API to fetch user roles and checks against device requirements.
- **Device Access Mapping**: Reads access requirements from `accessRequirementsForDevices.json`.
- **Health & Logging Endpoints**: For monitoring and debugging.
- **Environment Configuration**: Uses `.env` for sensitive credentials.
- **OpenAPI Spec**: API documented in `open-api-v1-spec.yaml`.

## API Endpoints

### Health & Logging
- `GET /health` — Service health check
- `GET /logger` — Logger endpoint (for debugging/logging info)

### Authorization API (v1)
- `POST /api/v1/authorize-for-door-access` — Authorize user for door access based on facial recognition data
- `GET /api/v1/requirements` — Get access requirements for devices

#### Authorization Flow
1. Receives user identification data from HikCentral alarm.
2. Queries WSO2 Identity Server (SCIM2) for user roles.
3. Validates user permissions against device requirements.
4. Returns authorization decision.

**Request Example:**
```json
{
  "id": "user123",
  "name": "John Doe",
  "description": "Facial recognition match from camera 1"
}
```

**Response Example:**
```json
{
  "id": "user123",
  "userName": "john.doe",
  "emails": [...],
  "groups": [...],
  "roles": [...],
  "authorized": true
}
```

## Getting Started

### Prerequisites
- Go 1.25.0 or later

### Environment Variables

Create a `.env` file in the project root:
```
IDP_ADDRESS=https://localhost:9443
IDP_USERNAME=admin
IDP_PASSWORD=admin
```

### Installation & Running

```bash
go mod tidy
go run main.go
# or
go build -o service main.go
./service
```

Service runs on port 5000 by default.

## Example Usage

**Authorize user for door access**
```bash
curl -X POST http://localhost:5000/api/v1/authorize-for-door-access \
  -H "Content-Type: application/json" \
  -d '{"id":"user123","name":"John Doe","description":"Facial recognition match from camera 1"}'
```

**Health check**
```bash
curl http://localhost:5000/health
```

## Adding New Features

1. **Create a Model**  
   Add a new file in `internal/models/` for your data structure.

2. **Create a Controller**  
   Add a new controller in `internal/controllers/v1_controllers/`.

3. **Create Routes**  
   Define new routes in `internal/routes/internal/v1_routes/`.

4. **Register Routes**  
   Update `internal/routes/routes.go` to register your new routes.

## Integration Points

- **HikCentral Professional**: Receives alarm notifications, processes user identification.
- **WSO2 Identity Server**: SCIM2 protocol for user data and role-based access control.

## Security Considerations

- TLS verification is disabled for development (enable for production).
- Basic authentication for SCIM2 API.
- Environment variables for credentials.
- No API key authentication for HikCentral webhooks (recommended for production).

## Next Steps for Production

- Enable TLS verification for SCIM2 API calls.
- Implement API key authentication for HikCentral webhooks.
- Add input validation and sanitization.
- Integrate door control hardware APIs.
- Add camera-to-door mapping configuration.
- Implement logging, monitoring, and audit logs.
- Add role-based authorization logic.
- Add unit and integration tests.
- Add Docker containerization.
- Add API documentation (Swagger/OpenAPI).
- Implement rate limiting and request throttling.
- Configure production-ready CORS policies.
