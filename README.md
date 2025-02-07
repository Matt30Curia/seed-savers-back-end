

# Seed Exchange Backend

This project provides the backend for a seed exchange platform, built using Go. It includes features for user authentication, data storage, and managing seed exchanges. The backend uses migrations for database management and supports JWT-based authentication for secure access.

## Installation

To get started with this project, you need to have Go installed on your system. If you don't have Go installed, you can download it from [here](https://golang.org/dl/).

Clone this repository and install the dependencies:

```bash
git clone <repository-url>
cd <repository-directory>
go mod tidy
```

---

## Environment Variables

This project requires several environment variables to be set in your environment or `.env` file:

| **Variable**             | **Description**                                              |
|--------------------------|--------------------------------------------------------------|
| `DB_PASSWORD`            | Database password.                                           |
| `DB_USER`                | Database username.                                           |
| `DB_HOST`                | Host of the database.                                        |
| `DB_PORT`                | Database port.                                               |
| `DB_NAME`                | Database name.                                               |
| `PUBLIC_HOST`            | The public URL of the application (e.g., `http://localhost`). |
| `JWT_SECRET`             | Secret key for JWT signing and validation.                   |
| `GOOGLE_CLIENT_ID`       | Google OAuth Client ID for authentication.                   |
| `GOOGLE_CLIENT_SECRET_ID`| Google OAuth Client Secret for authentication.               |
| `EMAIL`                  | Email address for sending notifications.                     |
| `EMAIL_PASSWORD`         | Email password for authentication.                           |
| `HOST_SMTP`              | SMTP server for sending emails.                              |
| `AES_KEY`                | Key for encryption (e.g., used for storing sensitive data).  |
| `TOKEN_EXPIRATION_HOUR`  | Token expiration time in hours for JWTs.                     |

You can configure these variables by setting them in a `.env` file or manually in your environment.

---

## Make Commands

This project uses `Make` for task automation. Below are some useful commands you can run:

- **Build the project:**
    ```bash
    make build
    ```
    Compiles the Go application into a binary at `bin/seed.exe`.

- **Run tests:**
    ```bash
    make test
    ```
    Runs the Go tests in the repository with verbose output.

- **Run the application:**
    ```bash
    make run
    ```
    Builds and runs the application.

- **Create a new migration:**
    ```bash
    make migration <migration-name>
    ```
    Generates a new SQL migration in the `cmd/migrate/migrations` directory.

---

## Database Migrations

The project uses migrations to manage the database schema.

- **Apply migrations up:**
    ```bash
    make migrate-up
    ```
    Applies all pending migrations to the database.

- **Apply migrations down:**
    ```bash
    make migrate-down
    ```
    Rolls back the most recent migration.

---

## Running the Application

Once the environment variables are set, and the database is configured, you can run the application.

1. Build the application:
    ```bash
    make build
    ```

2. Run the application:
    ```bash
    make run
    ```

This will start the backend server, and the application should be accessible based on the `PUBLIC_HOST` you configured in the environment variables.
