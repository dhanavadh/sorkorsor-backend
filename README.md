# Sorkorsor Backend

Backend API for **Sorkorsor**, a personal New Year messaging project that lets users create customized messages and attach images for their friends.

The service is written in Go with Gin, stores message data in PostgreSQL through GORM, and uploads images to Cloudflare R2 using its S3-compatible API.

## Features

- Create and retrieve personalized messages
- Generate short random IDs for message links
- Upload images directly through the API
- Generate presigned URLs for client-side uploads
- Automatically migrate the PostgreSQL message table on startup
- Health-check endpoint for deployments
- Multi-stage Docker build

## Tech Stack

- **Language:** Go 1.25
- **Web framework:** Gin
- **Database:** PostgreSQL
- **ORM:** GORM
- **Object storage:** Cloudflare R2
- **Storage SDK:** AWS SDK for Go v2 (S3-compatible API)
- **Containerization:** Docker

## Project Structure

```text
.
├── cmd/          # Application entry point and route registration
├── config/       # Environment-based configuration
├── database/     # PostgreSQL connection and migrations
├── handlers/     # HTTP request handlers
├── models/       # Database models
├── repository/   # Data-access layer
├── storage/      # Cloudflare R2 integration
├── utils/        # Password utilities
└── Dockerfile
```

## Getting Started

### Prerequisites

- Go 1.25 or later
- PostgreSQL
- A Cloudflare R2 bucket with public access configured

### 1. Clone the repository

```bash
git clone https://github.com/dhanavadh/sorkorsor-backend.git
cd sorkorsor-backend
```

### 2. Configure environment variables

Create a `.env` file in the project root:

```env
SERVER_PORT=8080
GIN_MODE=debug

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=sorkorsor
DB_SSLMODE=disable

R2_ACCOUNT_ID=your_cloudflare_account_id
R2_ACCESS_KEY=your_r2_access_key
R2_SECRET_ACCESS_KEY=your_r2_secret_access_key
R2_BUCKET_NAME=your_bucket_name
R2_PUBLIC_URL=https://your-public-r2-domain.example.com
```

`PORT` takes precedence over `SERVER_PORT` when both are set, which makes the service compatible with platforms that inject a port automatically.

### 3. Install dependencies and run

```bash
go mod download
go run ./cmd/main.go
```

The API starts at `http://localhost:8080` by default. GORM creates or updates the `messages` table automatically during startup.

## Running with Docker

Build the image:

```bash
docker build -t sorkorsor-backend .
```

Run the container with your environment file:

```bash
docker run --env-file .env -p 8080:8080 sorkorsor-backend
```

The container still requires access to the PostgreSQL database and Cloudflare R2 endpoints defined in `.env`. When PostgreSQL runs on the host, set `DB_HOST` to a hostname reachable from the container instead of `localhost`.

## API Reference

Base path: `/api/v1`

### Health check

```http
GET /health
```

Response:

```json
{
  "status": "ok"
}
```

### List messages

```http
GET /api/v1/messages
```

Returns all stored messages.

### Get a message

```http
GET /api/v1/message/:id
```

Returns the message associated with the supplied message ID.

### Create a message

```http
POST /api/v1/message
Content-Type: application/json
```

Example request:

```json
{
  "password": "12345",
  "recipient": "Friend's name",
  "profile_sticker": 1,
  "inventory_sticker": 2,
  "message": "Happy New Year!",
  "image_url": [
    "https://images.example.com/images/photo.jpg"
  ]
}
```

Example response:

```json
{
  "id": "generated_id"
}
```

### Upload a file through the API

```http
POST /api/v1/upload
Content-Type: multipart/form-data
```

Send the file using the form field `file`. A successful request returns its public URL:

```json
{
  "message": "upload successfully",
  "url": "https://images.example.com/uploaded-file.jpg"
}
```

### Generate presigned upload URLs

```http
POST /api/v1/upload/presign
Content-Type: application/json
```

Example request:

```json
{
  "files": [
    {
      "content_type": "image/jpeg",
      "extension": ".jpg"
    },
    {
      "content_type": "image/png",
      "extension": ".png"
    }
  ]
}
```

Each generated upload URL is valid for 15 minutes. Example response:

```json
{
  "urls": [
    {
      "upload_url": "https://presigned-upload-url.example.com",
      "key": "images/generated-key.jpg",
      "public_url": "https://images.example.com/images/generated-key.jpg"
    }
  ]
}
```

Upload the corresponding file to `upload_url` with an HTTP `PUT` request and the same `Content-Type` used when requesting the URL.

## CORS

The current configuration allows requests from the local frontend at `http://localhost:5173`. Update the allowed origins in `cmd/main.go` before deploying with a different frontend URL.

## Future Improvements

- Add request validation and consistent error responses
- Hash message passwords before storing them
- Prevent sensitive fields from being returned by read endpoints
- Add automated unit and integration tests
- Add OpenAPI/Swagger documentation
- Configure CORS through environment variables
- Add CI checks for formatting, tests, and builds

## License

This project does not currently include a license.
