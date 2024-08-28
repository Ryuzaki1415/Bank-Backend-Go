# Bank Backend Server

This is a backend server for a banking application written in Go. It provides basic banking operations and uses JWT for authentication.

## Key Features

### 1. User Account Management
- **Create Account**: Users can open new bank accounts.
- **Delete Account**: Existing accounts can be closed upon request.
- **Search Account**: Accounts can be retrieved using their unique account numbers.

### 2. Authentication System
- **User Login**: Secure login functionality for account holders.
- **JWT Authentication**: All sensitive operations are protected using JSON Web Tokens (JWT).

### 3. Technology Stack
- **Go (Golang)**: The entire backend is written in Go, known for its performance and concurrency support.
- **No External Libraries**: The project relies solely on Go's standard library, ensuring minimal dependencies and maximum control.

### 4. Database
- **PostgreSQL**: Utilizes PostgreSQL for robust and reliable data storage.
- **Docker Integration**: The PostgreSQL database is hosted in a Docker container, enabling easy setup and deployment.

### 5. Security
- **JWT Implementation**: Custom JWT implementation for secure authentication without relying on external libraries.
- **Protected Routes**: All sensitive operations require valid JWT tokens.

### 6. Scalability and Performance
- **Efficient Data Handling**: Optimized database queries and connection management.
- **Concurrent Request Handling**: Leverages Go's goroutines for handling multiple requests simultaneously.

## Prerequisites

- Go (version 1.19 or later)
- Gorilla/Mux for routing
- Docker
- PostgreSQL
- Postman/Thunderclient

## Setup

1. Clone the repository:
   ```
   git clone https://github.com/Ryuzaki1415/Bank-backend-Go.git
   cd Bank-backend-Go
   ```

2. Start the PostgreSQL database using Docker:
   ```
   docker run --name bank-postgres -e POSTGRES_PASSWORD=yourpassword -p 5432:5432 -d postgres
   ```

3. Install dependencies:
   ```
   go mod download
   ```

4. Set up environment variables:
   ```
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=postgres
   export DB_PASSWORD=yourpassword
   export DB_NAME=bankdb
   export JWT_SECRET=your_jwt_secret_key
   ```

5. Run the server:
   ```
   go run . //to run all modules.
   ```

## API Endpoints

- `POST /account`: Create a new user account with body elements {firstName, lastName, password}
- `POST /login`: Authenticate and receive a JWT token
- `DELETE /account`: Delete user account (requires authentication)
- `GET /account/id`: Search for an account by account number (requires authentication)

## Authentication

All protected endpoints require a valid JWT token in the Authorization header.
As of now, only the endpoint GET /account/id is protected.