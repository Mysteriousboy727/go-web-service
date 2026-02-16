# Go Industry Server

A RESTful API server built with Go, Gin, and GORM following clean architecture principles.

## Project Structure

```
go-industry-server/
ÃÄÄ cmd/server/          # Application entry point
ÃÄÄ internal/            # Private application code
³   ÃÄÄ handler/         # HTTP handlers
³   ÃÄÄ service/         # Business logic
³   ÃÄÄ repository/      # Data access layer
³   ÃÄÄ middleware/      # HTTP middleware
³   ÀÄÄ model/           # Domain models
ÃÄÄ pkg/                 # Public reusable packages
ÃÄÄ configs/             # Configuration files
ÀÄÄ .env                 # Environment variables
```

## Setup

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Set up your database and update `.env` file

3. Run the server:
   ```bash
   go run cmd/server/main.go
   ```

## API Endpoints

- `GET /health` - Health check
- `POST /api/v1/users` - Create user
- `GET /api/v1/users` - Get all users
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user
