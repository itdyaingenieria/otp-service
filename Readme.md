
# OTP Service

A robust, production-ready One-Time Password (OTP) microservice built with Go. This service provides secure OTP generation, validation, and delivery, designed for integration with authentication and verification workflows (e.g., 2FA, user onboarding, transaction confirmation).

## Features

- Generate numeric OTP codes for any channel (SMS, email, etc.)
- Validate OTPs with expiration, attempt limits, and single-use enforcement
- PostgreSQL persistence for auditability and reliability
- Modular architecture (ports & adapters, clean architecture)
- Dockerized for easy deployment
- Ready for extension to real notification providers (currently logs to console)

## Architecture Overview

- **cmd/api/main.go**: Application entrypoint, HTTP server setup
- **internal/domain**: Core business logic (OTP entity)
- **internal/usecase**: Application use cases (generate, validate OTP)
- **internal/adapter**: Infrastructure (clock, codegen, notifier, repository)
- **internal/httpapi**: HTTP handlers and routing
- **internal/config**: Configuration loading
- **migrations/**: Database schema (PostgreSQL)

## Quick Start

### Prerequisites
- Docker & Docker Compose installed

### 1. Clone the repository
```bash
git clone <your-repo-url>
cd otp-service
```

### 2. Configure Environment
Edit the `.env` file as needed. Example:
```env
PORT=8098
DB_HOST=db
DB_PORT=5432
DB_USER=otp_user
DB_PASSWORD=otp_pass
DB_NAME=otp_db
```

### 3. Start the Service
```bash
docker compose up --build
```
- The OTP service will be available at `http://localhost:8098`
- PostgreSQL will be available at `localhost:5432`

## API Endpoints

### Generate OTP
- **POST** `/api/v1/otp`
- **Body:**
	```json
	{
		"tenant_id": "string",
		"channel": "sms|email",
		"destination": "string"
	}
	```
- **Response:**
	```json
	{
		"id": "string",
		"expires_at": "RFC3339 timestamp"
	}
	```

### Validate OTP
- **POST** `/api/v1/otp/validate`
- **Body:**
	```json
	{
		"tenant_id": "string",
		"id": "string",
		"code": "string"
	}
	```
- **Response:**
	- 200 OK if valid
	- 400/401/404 with error message if invalid/expired/used

## Database
- Uses PostgreSQL
- Schema is auto-applied from `migrations/001_create_otps.sql`

## Development

### Run Locally (without Docker)
1. Start PostgreSQL manually and apply migrations
2. Set environment variables or edit `.env`
3. Build and run:
	 ```bash
	 go build -o otp-service ./cmd/api/main.go
	 ./otp-service
	 ```

### Run Tests
```bash
go test ./...
```

## Extending
- Implement real notification providers in `internal/adapter/notifier/`
- Add more channels or OTP formats in `internal/adapter/codegen/`

## License
MIT

## Authors
- Itdyaingenieria - diegoyamaa@gmail.com

---

For questions or contributions, please open an issue or pull request.