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
    goose -dir .\internal\database\migrations postgres "postgres://postgres:PASSWORD@localhost:5432/postgres" up
    ```

6. **Run the application:**

    ```sh
    cd .\cmd\server
    go run .
    ```

    The server will start running on `http://localhost:8080`.

## API Endpoints

### Authentication

- **POST /register** - Register a new user.
    
    *Request*

    ```js
    {
    "login": string,
    "password": string
    }
    ```

    *Response*

    ```js
    {
    "message": "User registered successfully",
    "token": string
    "user_id": int
    }
    ```


- **POST /login** - Login an existing user.

    *Request*

    ```js
    {
    "login": string,
    "password": string
    }
    ```

    *Response*

    ```js
    {
    "message": "Login successful",
    "token": string
    "user_id": int
    }
    ```

### Results

- **POST /results** - Create a new result (requires authentication).

    *Request*

    **header**
    | Key       | Value       |
    |-----------|-------------|
    | token     | string      |

    
    ```js
    {
    "game_mode": string,
    "duration": time.Time,
    "mistakes": int,
    "accuracy": float,
    "count_words": int,
    "wpm": float,
    "cpm": float
    }
    ```

    *Response*

    ```js
    {
    "message": "Result created successfully"
    }
    ```

- **GET /results/** - Get all results for a user.

    *Request*

    **Query Params**
    | Key       | Value       |
    |-----------|-------------|
    | user_id   | int         |

    *Response*

    ```js
    [
    {
        "id": int,
        "user_id": int,
        "game_mode": string,
        "start_time": time.Time,
        "duration": time.Time,
        "mistakes": int,
        "accuracy": float,
        "count_words": int,
        "wpm": float,
        "cpm": float
    }
    ]
    ```


- **GET /results/:id** - Get a specific result by ID.

    *Response*

    ```js
    [
    {
        "id": int,
        "user_id": int,
        "game_mode": string,
        "start_time": time.Time,
        "duration": time.Time,
        "mistakes": int,
        "accuracy": float,
        "count_words": int,
        "wpm": float,
        "cpm": float
    }
    ]
    ```

### Records

- **GET /records/:gamemode** - Get global records for a specific game mode.

    *Response*

    ```js
    [
    {
        "id": int,
        "user_id": int,
        "game_mode": string,
        "start_time": time.Time,
        "duration": time.Time,
        "mistakes": int,
        "accuracy": float,
        "count_words": int,
        "wpm": float,
        "cpm": float
    }
    ]
    ```

- **GET /records** - Get user-specific records.

    *Request*

    **Query Params**
    | Key       | Value       |
    |-----------|-------------|
    | user_id   | int         |


    *Response*

    ```js
    [
    {
        "id": int,
        "user_id": int,
        "game_mode": string,
        "start_time": time.Time,
        "duration": time.Time,
        "mistakes": int,
        "accuracy": float,
        "count_words": int,
        "wpm": float,
        "cpm": float
    }
    ]
    ```
