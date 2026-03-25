# Go Gully Backend

A Go-GIN REST API backend with Google OAuth authentication and Gully integration.


## Dependent Repos
- [Gully Core Agent](https://github.com/Astrasv/the-gully)

## Prerequisites

- Go 1.25+
- PostgreSQL
- Google OAuth credentials

## Installation

1. Clone the repository:
```bash
git clone https://github.com/Astrasv/go-gully-backend
cd go-gully-backend
```

2. Install dependencies:
```bash
go mod download
```

3. Copy the example environment file:
```bash
cp .env.example .env
```

4. Configure your `.env` file:
```env
PORT=8080
DB="host=localhost user=postgres password=password dbname=gully_backend port=5432 sslmode=disable"

GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_CALLBACK_URL=http://localhost:8080/auth/google/callback
GOOGLE_REDIRECT_URL=http://localhost:8080
SESSION_SECRET=your_session_secret

SQL_AGENT_URL=http://localhost:8000

FRONTEND_URL=http://localhost:3000
```

5. Set up the database:
```bash
# Create PostgreSQL database if not created
createdb gully_backend
```

## Running the Application

### Local Development
```bash
go run main.go
```

### Using Docker
```bash
docker-compose up --build
```


## License

MIT