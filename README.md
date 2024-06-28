# Agency-Type Backend

## Overview

This repository contains the backend code for the Agency-Type application. It provides a RESTful API to handle user authentication, results management, and records retrieval. The backend is built using Go, with the Echo framework for handling HTTP requests, and PostgreSQL for database operations.

## Features

- **User Authentication**: Register and login endpoints with JWT token generation.
- **Results Management**: Create, retrieve, and manage user results.
- **Records Retrieval**: Fetch global and user-specific records based on game modes.

## Technologies Used

- **Go**: The primary programming language.
- **Echo**: A high performance, extensible, and minimalist web framework for Go.
- **PostgreSQL**: A powerful, open source object-relational database system.
- **JWT**: JSON Web Tokens for secure authentication.
- **bcrypt**: For password hashing.

## Prerequisites

- Go (version 1.16 or higher)
- PostgreSQL (version 10 or higher)
- Git

## Installation

1. **Clone the repository:**

    ```sh
    git clone https://github.com/your-repo/agency-type-backend.git
    cd agency-type-backend
    ```

2. **Set up the environment variables:**

    Create a `.env` file in the root directory and add the following:

    ```env
    DATABASE_URL=postgres://postgres:PASSWORD@localhost:5432/postgres?sslmode=disable
    SECRET_TOKEN="TOKEN"
    ```

3. **Install dependencies:**

    ```sh
    go mod download
    ```

4. **Install Goose for database migrations:**

    ```sh
    go install github.com/pressly/goose/v3/cmd/goose@latest
    ```

5. **Run the database migrations:**

    ```sh
    goose -dir \internal\pkg\database\migrations postgres "postgres://postgres:123@localhost:5432/postgres" up
    ```

6. **Run the application:**

    ```sh
    cd .\cmd\server
    go run main.go
    ```

    The server will start running on `http://localhost:8080`.

## API Endpoints

### Authentication

- **POST /register** - Register a new user.
- **POST /login** - Login an existing user.

### Results

- **POST /results** - Create a new result (requires authentication).
- **GET /results** - Get all results for a user.
- **GET /results/:id** - Get a specific result by ID.

### Records

- **GET /records/:gamemode** - Get global records for a specific game mode.
- **GET /records** - Get user-specific records.

