# people-enrichment-service

A service for enriching personal information (age, gender, nationality) by full name using public APIs and saving the results to PostgreSQL.

## Features

- Add a new person (enrichment via public APIs)
- Get a list of people with filters and pagination
- Get person info by ID
- Update person information
- Delete person by ID
- Logging (zap)
- Swagger documentation (`docs/swagger.yml`)
- Configuration via `.env`

## Public APIs Used

- [Agify.io (age)](https://api.agify.io)
- [Genderize.io (gender)](https://api.genderize.io)
- [Nationalize.io (nationality)](https://api.nationalize.io)

## Quick Start

### 1. Clone the repository

```sh
git clone https://github.com/your-username/people-enrichment-service.git
cd people-enrichment-service
```

### 2. Fill in the `.env` file

Example:
```
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=people
```

### 3. Start the service and database with Docker

```sh
make up
```

The service will be available at [http://localhost:8080](http://localhost:8080)

### 4. Stop the service

```sh
make down
```

### 5. Local run (without Docker)

- Make sure PostgreSQL is running and matches your `.env` settings
- Apply migrations (e.g., via `psql` or a migration tool)
- Run the application:

```sh
make run
```

### 6. Makefile commands help

To see all available Makefile commands:

```sh
make help
```

You will see a list of commands for managing the service, containers, and building the project.

## API Documentation

Swagger spec is available at [`docs/swagger.yml`](docs/swagger.yml).  
You can visualize it using [Swagger Editor](https://editor.swagger.io/).

## How to view Swagger

1. Start the service:
   ```sh
   make up
   ```
2. Open [http://localhost:8081](http://localhost:8081) in your browser
3. Enter `/swagger.yml` in the search field and press Enter — your API documentation will open.

## Example Requests

**Add a person:**
```http
POST /person
Content-Type: application/json

{
  "name": "Dmitriy",
  "surname": "Ushakov",
  "patronymic": "Vasilevich"
}
```

**Get a list:**
```http
GET /person?limit=10&offset=0&name=Dmitriy
```

## Tests

To run unit tests:
```sh
go test ./... -v
```

## Project Structure

```
.
├── cmd/app/                # Application entry point
├── internal/
│   ├── handler/            # HTTP handlers (controllers)
│   ├── service/            # Business logic
│   ├── repository/         # Database access
│   ├── models/             # Data models
│   ├── config/             # Configuration loader
│   └── logger/             # Logger initialization and configuration
├── migrations/             # SQL migrations for DB
├── docs/                   # Swagger and other docs
├── utils/                  # Utility functions
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── README.md
```

---

**Author:**  
adal4ik