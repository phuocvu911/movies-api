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
If database got stale between the run (you want to GET the entity that is already got force DELETE), restart the server with:

```bash
rm data.db && go run main.go
```
## How It Works
The server listens on the port `:8080` by default. On first run it creates `data.db`, migrates the schema and seeds it with sample data. Seeding is skipped automatically if the database already contains movies.

## API reference

All three entities (`genres`, `movies`, `actors`) support the same CRUD pattern:

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/{entity}` | Create (201) |
| GET | `/api/{entity}` | List all, paginated (200) |
| GET | `/api/{entity}/{id}` | Get one by id (200 / 404) |
| PATCH | `/api/{entity}/{id}` | Partial update (200) |
| DELETE | `/api/{entity}/{id}` | Delete (400/204); add `?force=true` to also remove relationships |

Relationship and filter endpoints:

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/movies?genre={genreId}` | Movies in a genre |
| GET | `/api/movies?year={releaseYear}` | Movies from a year |
| GET | `/api/movies?actor={actorId}` | Movies an actor starred in |
| GET | `/api/movies/{id}/actors` | All actors in a movie |
| GET | `/api/movies/search?title={text}` | Case-insensitive partial movie title search |
| GET | `/api/actors?name={text}` | Case-insensitive partial actor name search |
| GET | `/api/actors/{id}/movies` | Movies an actor appeared in |
| GET | `/api/genres/{id}/movies` | Movies having a genre |

Filters for `movies` can be combined, e.g. `/api/movies?genre=1&year=1999`.

### Deletion and relationships

By default, deleting an entity that still has relationships fails with `400`:

```
DELETE /api/genres/1
→ 400 "Unable to delete genre 'Action' as it is associated with 12 movies"

DELETE /api/actors/6
→ 400 "Unable to delete actor 'Tom Hanks' as he/she is associated with 2 movies"
```

Force deletion removes the entity and it's relationships:

```
DELETE /api/genres/1?force=true
→ 204 No Content
```

### Error handling

| Status | Meaning |
|--------|---------|
| 400 | Validation failure, malformed body/params, or blocked deletion |
| 404 | Entity (or referenced filter entity) does not exist |
| 409 | Conflict with the table (when user try to create duplicate `genre` or `actor`- with exact name and birthdate) |
| 500 | Unexpected server error |

Validation covers required fields, release year max 2027, positive duration, ISO 8601 birth dates, birthdate can not be in the future, and positivity of every referenced genre/actor/movie id.

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
- Implements per-IP rate limiting using [token bucket algorithm](https://www.geeksforgeeks.org/system-design/rate-limiting-algorithms-system-design/) with the `golang.org/x/time/rate` package.
- Limit: 100 requests per second with a burst capacity of 150 requests (you can adjust these values in the `rateLimit` middleware function to see the effect).
- Returns `429 Too Many Requests` when rate limit is exceeded
- Prevents abuse and ensures fair API usage
- There is also a go routine run in the background to clean up the rate limiters for IPs that haven't made requests in the last 3 minutes, preventing memory bloat.

The flow of requests when hitting the API is as follows:
```
request → Logging → RateLimit → RequireJSON →  handler
```
