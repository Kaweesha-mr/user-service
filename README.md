# User Service Documentation

This document provides a comprehensive guide to setting up and using the User Service. Follow these steps to get the service running and understand its key components.

## Table of Contents
1. [Overview](#overview)
2. [Features](#features)
3. [Prerequisites](#prerequisites)
4. [Setup Instructions](#setup-instructions)
   - [Using Docker Compose](#using-docker-compose)
   - [Environment Variables](#environment-variables)
5. [Key Components](#key-components)
6. [Dependencies](#dependencies)
7. [Future Enhancements](#future-enhancements)

---

### Overview
The User Service is a backend service designed for managing users. It supports operations such as creating, updating, retrieving, and deleting user records.

**Caching**: Redis is used to cache verified tokens and frequently accessed user data, improving performance.

**Token Validation**: Token validation requires a separate Authentication Service. The User Service will not work correctly without the Authentication Service being set up.

---

### Features
- User CRUD operations
- Redis caching for performance optimization
- Health check endpoint to monitor service status
- JWT Middleware for secure access
- Dependency injection for modular and testable code

---

### Prerequisites
1. Docker and Docker Compose installed on your system.
2. PostgreSQL and Redis configured via Docker Compose.
3. An Authentication Service for token validation.

---

### Setup Instructions

#### Using Docker Compose
The service relies on PostgreSQL and Redis. Use the provided `docker-compose.yml` file to set up the required containers.

1. **Download or create the `docker-compose.yml` file** with the following content:
   ```yaml
   services:
     postgres:
       image: postgres:latest
       container_name: postgres_container_users_v2
       environment:
         POSTGRES_USER: test
         POSTGRES_PASSWORD: 1234
         POSTGRES_DB: users
       ports:
         - "5433:5432"
       volumes:
         - postgres_data:/var/lib/postgresql/data
       networks:
         - mynetwork

     pgadmin:
       image: dpage/pgadmin4
       container_name: pgadmin_container_users_v2
       environment:
         PGADMIN_DEFAULT_EMAIL: admin@admin.com
         PGADMIN_DEFAULT_PASSWORD: admin
       ports:
         - "5051:80"
       depends_on:
         - postgres
       networks:
         - mynetwork

     redis:
       image: redis:latest
       container_name: redis_container_users_v2
       ports:
         - "6379:6379"
       networks:
         - mynetwork

   volumes:
     postgres_data:

   networks:
     mynetwork:
       driver: bridge
   ```

2. **Start the services**:
   ```bash
   docker-compose up -d
   ```

3. Verify that the PostgreSQL, PgAdmin, and Redis containers are running.

#### Environment Variables
Configure the service using the `.env` file. Below is the required configuration:
```env
DB_DSN=host=localhost user=test password=1234 dbname=users port=5433 sslmode=disable TimeZone=Asia/Colombo
```

---

### Key Components

#### Router
The `SetUpRouter` function initializes the Gin router and sets up routes, middleware, and controllers. It connects to both PostgreSQL and Redis.

Example routes:
- **Health Check**: `GET /`
- **User Operations**:
  - `GET /v1/users` - Get all users
  - `GET /v1/user/:id` - Get user by ID
  - `POST /v1/users` - Create a new user
  - `PUT /v1/update` - Update a user
  - `DELETE /v1/delete/:id` - Delete a user

#### Caching with Redis
- Used to cache verified tokens and user results.
- Improves the performance of frequent operations.

#### Token Validation
- Requires a separate Authentication Service.
- Without the Authentication Service, token validation will fail.

---

### Dependencies
- **Gin Framework**: For building HTTP web services.
- **Redis**: For caching.
- **PostgreSQL**: Primary database.
- **Docker Compose**: For managing dependencies and infrastructure.

---

### Future Enhancements
- Addition of more endpoints for advanced user management.
- Improved role-based access control.
- Enhanced logging and monitoring.

---

### Running the Application
Once the database and caching layers are up, start the User Service:
```bash
go run main.go
```
Access the health check endpoint to verify the service is running:
```bash
curl http://localhost:8080/
```
Expected response:
```json
{
  "status": "Service is up and running"
}
```

