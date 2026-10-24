# DroneMEQ Violation Management Backend

Backend service for managing drone violations and related enforcement actions. Built with Go, Gin, Gorm, and follows Clean Architecture principles.

## Table of Contents
- [Overview](#overview)
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
  - [Using Docker](#using-docker)
  - [Running Locally](#running-locally)
- [Environment Variables](#environment-variables)
- [Database Migration & Seeding](#database-migration--seeding)
- [API Documentation](#api-documentation)
- [ERD Diagram](#erd-diagram)
- [License](#license)

## Overview
This service provides RESTful APIs to:
- Record drone violation incidents (location, type, severity, reporter, etc.)
- Manage violation actions (warnings, fines, etc.)
- Attach media and notes to violations
- Configure offices and users (via constants and helper packages)
- Send notifications (email via SMTP)

## Architecture
The codebase is organized into three main layers:

### Data Layer (`data-layer`)
Handles persistence and database interaction.
- `config`: Database initialization, migration, seeding, and third-party configs (e.g., Cloudinary).
- `model`: Gorm models representing tables (Violation, ViolationAction, ViolationMedia, ViolationNote, etc.).
- `repository`: Data access interfaces and implementations for each model.
- `seeder`: Seed data (JSON) for initial setup.

### Domain Layer (`domain-layer`)
Contains core business logic.
- `service`: Business logic layer that orchestrates repository calls and enforces rules.
- `middleware`: HTTP middleware (authentication, validation, logging, etc.).
- `client`: External service clients (if any).

### Presentation Layer (`presentation-layer`)
Handles HTTP transport.
- `controller`: HTTP handlers that validate requests, call services, and format responses.
- `routes`: Gin route definitions mapping endpoints to controllers.

## API Endpoints
All API endpoints are prefixed with `/api/violations` and require authentication via JWT (handled by middleware).

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/` | Create a new violation record |
| GET | `/` | Retrieve all violations (supports pagination via query params) |
| GET | `/{Pkid}` | Get a single violation by its primary key |
| GET | `/filter` | Filter violations by query parameters (e.g., type, status, date range) |
| PUT | `/{Pkid}` | Update an existing violation |
| PUT | `/{Pkid}/soft-delete` | Soft‑delete a violation (sets `is_deleted = true`) |
| PUT | `/{Pkid}/restore` | Restore a soft‑deleted violation |
| DELETE | `/{Pkid}` | Hard‑delete a violation from the database |
| POST | `/{Pkid}/media` | Upload media (photo/video) associated with a violation |
| POST | `/{Pkid}/notes` | Add a note to a violation |
| GET | `/{Pkid}/notes` | Retrieve all notes for a violation |
| GET | `{Pkid}/notes/filternote` | Retrieve notes filtered by type/category |
| PUT | `/{Pkid}/notes/{notePkid}` | Update a specific note |
| PUT | `/{Pkid}/notes/{notePkid}/soft-delete` | Soft‑delete a note |
| PUT | `/{Pkid}/notes/{notePkid}/restore` | Restore a soft‑deleted note |
| DELETE | `/{Pkid}/notes/{notePkid}` | Hard‑delete a note |
| POST | `/{Pkid}/actions` | Record an enforcement action (warning, fine, etc.) for a violation |


## Prerequisites
- Go version >= 1.21
- Docker & Docker Compose (optional, for containerized setup)
- PostgreSQL >= 8.0 (if running locally without Docker)

## Getting Started

### Using Docker
1. Copy the example environment file and configure it:
   ```bash
   cp .env.example .env
   ```
2. Edit `.env` with your credentials (or leave defaults for local dev).
3. Build and start the services:
   ```bash
   docker compose up -d --build
   ```
   This will start the Go app (port 8888) and PostgreSQL (port 5432, mapped from `${DB_PORT}`).

### Running Locally
1. Ensure PostgreSQL is running and accessible.
2. Copy and configure `.env`:
   ```bash
   cp .env.example .env
   ```
   Adjust `DB_HOST`, `DB_USER`, `DB_PASS`, `DB_NAME`, `DB_PORT` as needed.
3. Install dependencies:
   ```bash
   go mod tidy
   ```
4. Run the application:
   ```bash
   go run main.go
   ```
   The server will start on `:8888` (or port from `GOLANG_PORT`).

## Environment Variables
See `.env.example` for reference. Key sections:
- **Database**: Connection details for PostgreSQL.
- **Server**: `NGINX_PORT` (if using reverse proxy) and `GOLANG_PORT`.
- **SMTP**: Settings for outgoing email.
- **Cloudinary**: Configuration for media uploads.
- **User Management**: URL for external user service (if integrated).

## Database Migration & Seeding
The binary accepts flags to manage the database:
- **Migrate**: Apply pending migrations.
  ```bash
  go run main.go --migrate
  ```
- **Seed**: Load initial data from `data-layer/seeder/`.
  ```bash
  go run main.go --seed
  ```
- **Both**: 
  ```bash
  go run main.go --migrate --seed
  ```

> Note: Migration and seeder logic resides in `data-layer/config/migrate.go` and `seed.go`.

## API Documentation
API endpoints are defined in `presentation-layer/routes/` and implemented in controllers.
For detailed API reference, see the Postman collection (to be added) or generate Swagger docs via annotations (future work).

## ERD Diagram
Entity-Relationship diagram illustrating table relationships:
![ERD Diagram](https://dbdiagram.io/d/Violation-Management-66f0e910a0828f8aa6b20b65)

Alternatively, view the live diagram at: https://dbdiagram.io/d/Violation-Management-66f0e910a0828f8aa6b20b65

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.