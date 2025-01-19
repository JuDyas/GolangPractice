# Pastebin Project
## Description
Pastebin is a web application for storing, managing, and sharing text data. Users can create, delete, and view entries (pastes). Administrators have additional privileges, such as access to all entries and filtering capabilities.

---
## Technologies

#### Backend
- **Programming Language**: [Go](https://go.dev/)
- **HTTP Framework**: [Gin](https://github.com/gin-gonic/gin)
- **Database**: [MongoDB](https://github.com/mongodb/mongo-go-driver)
- **Authentication**: [JWT](https://en.wikipedia.org/wiki/JSON_Web_Token)

#### Frontend
- **In feature updates**

## Features

### Users
- **Registration**
- **Authentication** via JWT
- **Creating entries (pastes)** with parameters such as password protection, access only for specified email, IP, only for authorized users
- **Viewing entries**

### Administrators
- **Soft delete entries**
- **Retrieve a list of entries with filtering and pagination**
- **Filter by columns and sort entries**

---

## Installation and Setup

### 1. Clone the Repository

### 2. Install Dependencies
```bash
go mod tidy
```

### 3. Configure the Application
1. Open the `config.json` file.
2. Set the MongoDB connection parameters and lockal port:
    ```json
    {
        "port": ":8080",
        "uri": "mongodb://mongo:27017",
    }
    ```

### 4. Run the Application
```bash
make run
```
The application will be available at `http://localhost:8080`.

---

## API Endpoints

### Authentication
- **POST /v1/auth/register** — Register a user.
- **POST /v1/auth/login** — Authenticate a user.

### Entries (Pastes)
- **POST /v1/pastes** — Create an entry.
- **GET /v1/:id** — Retrieve an entry by ID.

### Administration
- **DELETE /v1/pastes/:id** — Soft delete an entry.
- **POST /v1/pastes/admin** — Retrieve a list of entries with filtering and pagination.

---

## Testing

**In feature updates**

---

## Docker

### Build and run the Docker Image
```bash
make up
```
