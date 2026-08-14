Phuoc: actor, genre, basic search

Hang: movies, pagination

# movies-api

A REST API for managing a movie database, built for the local film society. It stores movies, genres and actors in SQLite with full many-to-many relationships, and exposes CRUD plus filtering, search and pagination over plain HTTP/JSON.

## Setup and installation

### Requirement
- Go 1.26+
- C complier (if you are using Windows)
- Postman(for testing the API server)

### Run
```bash
git clone https://gitea.kood.tech/hoangphuocvu/movies-api
cd movies-api
go run main.go
```

## How It Works
...

## Extras
### Pagination
...

### Basic Search

The API supports searching movies by title using a query parameter (matches partial titles, case-insensitive).
If the `title` query parameter is not provided or is empty, the API will return a `400 Bad Request` error since we already have a `GET /api/movies` endpoint that returns all movies. Using `search` implies you're searching for something specific, so an empty query is considered invalid.

**Endpoint:**
```
GET /api/movies/search?title={search_term}
```

**Usage:**
```bash
curl "http://localhost:8080/api/movies/search?title=iNcep"
```

**Response:**

```json
[
    {
        "id": 2,
        "title": "Inception",
        "release_year": 2010,
        "duration": 148
    }
]
```

### Middlewares

The API server includes three middlewares for request handling and protection:

**1. Logging Middleware**
- Logs each request's HTTP method, URL path, and execution duration.
- Provides visibility into API usage and performance right in terminal.

**2. JSON Content-Type Validation**
- Enforces `Content-Type: application/json` for POST and PATCH requests.
- Returns `415 Unsupported Media Type` error if the header is missing or incorrect.
- Ensures consistent request/response handling.

**3. Rate Limiting**
- Implements per-IP rate limiting using [token bucket algorithm](https://www.geeksforgeeks.org/system-design/rate-limiting-algorithms-system-design/)
- Limit: 100 requests per second with a burst capacity of 150 requests (you can adjust these values in the `rateLimit` middleware function to see the effect).
- Returns `429 Too Many Requests` when rate limit is exceeded
- Prevents abuse and ensures fair API usage

The flow of requests when hitting the API is as follows:
```
request → Logging → RateLimit → RequireJSON →  handler
```
