# Visit Tracker API

A RESTful API for caregiver visit tracking and Electronic Visit Verification (EVV) compliance, built with Go and Gin framework.

## Features

- 📅 Schedule management for caregiver visits
- 🕒 Start/end visit logging with timestamps and geolocation
- ✅ Task tracking and completion status
- 📊 Dashboard statistics and reporting
- 🗄️ SQLite database with automatic schema setup and sample data

## Tech Stack

- **Backend:** Go 1.21+ with Gin framework
- **Database:** SQLite
- **Dependencies:**
  - `github.com/gin-gonic/gin` - Web framework
  - `github.com/gin-contrib/cors` - CORS middleware
  - `github.com/mattn/go-sqlite3` - SQLite driver

## Quick Start

### Prerequisites

- Go 1.21 or higher
- CGO enabled (required for SQLite)

### Installation

1. Clone the repository and navigate to the server directory:
```bash
cd server
```

2. Initialize Go module and install dependencies:
```bash
go mod tidy
```

3. Run the application:
```bash
go run main.go
```

The server will start on port 8080 by default. You can change the port by setting the `PORT` environment variable.

### Database

The application automatically creates a SQLite database file (`visits.db`) with sample data on first run. The database includes:
- 4 sample schedules (including today's and yesterday's)
- 5 tasks per schedule
- Visit tracking records

## API Endpoints

### Health Check
- `GET /health` - Server health status

### Schedule Management
- `GET /api/v1/schedules` - Get all schedules
- `GET /api/v1/schedules/today` - Get today's schedules
- `GET /api/v1/schedules/:id` - Get schedule details with tasks and visit info
- `GET /api/v1/schedules/:id/tasks` - Get tasks for a specific schedule

### Visit Tracking
- `POST /api/v1/schedules/:id/start` - Start a visit
- `POST /api/v1/schedules/:id/end` - End a visit

### Task Management
- `POST /api/v1/tasks/:taskId/update` - Update task status

### Statistics
- `GET /api/v1/stats` - Get dashboard statistics

## API Usage Examples

### Start a Visit
```bash
curl -X POST http://localhost:8080/api/v1/schedules/1/start \
  -H "Content-Type: application/json" \
  -d '{"latitude": 40.7128, "longitude": -74.0060}'
```

### End a Visit
```bash
curl -X POST http://localhost:8080/api/v1/schedules/1/end \
  -H "Content-Type: application/json" \
  -d '{"latitude": 40.7128, "longitude": -74.0060}'
```

### Update Task Status
```bash
# Mark as completed
curl -X POST http://localhost:8080/api/v1/tasks/1/update \
  -H "Content-Type: application/json" \
  -d '{"status": "completed"}'

# Mark as not completed with reason
curl -X POST http://localhost:8080/api/v1/tasks/2/update \
  -H "Content-Type: application/json" \
  -d '{"status": "not_completed", "reason": "Client was not available"}'
```

### Get Statistics
```bash
curl http://localhost:8080/api/v1/stats
```

## Data Models

### Schedule
- **id**: Unique identifier
- **client_name**: Name of the client
- **shift_start**: Start time of the shift
- **shift_end**: End time of the shift
- **location**: Address of the visit
- **status**: `upcoming`, `in_progress`, `completed`, `missed`

### Task
- **id**: Unique identifier
- **schedule_id**: Associated schedule
- **description**: Task description
- **status**: `pending`, `completed`, `not_completed`
- **reason**: Required when status is `not_completed`

### Visit
- **id**: Unique identifier
- **schedule_id**: Associated schedule
- **start_time**: Timestamp when visit started
- **end_time**: Timestamp when visit ended
- **start_lat/start_lng**: GPS coordinates at start
- **end_lat/end_lng**: GPS coordinates at end

## Business Logic

1. **Visit Flow**: 
   - Schedule starts with `upcoming` status
   - Start visit changes status to `in_progress`
   - End visit changes status to `completed`

2. **Task Management**:
   - Tasks can only be updated when visit is `in_progress`
   - Reason is required when marking task as `not_completed`

3. **Geolocation**:
   - GPS coordinates are required for both start and end visits
   - Coordinates are stored for compliance tracking

## Development

### Environment Variables
- `PORT`: Server port (default: 8080)
- `GIN_MODE`: Gin framework mode (`debug`, `release`, `test`)

### Database Reset
To reset the database with fresh sample data:
```bash
rm visits.db
go run main.go
```

## Error Handling

The API returns appropriate HTTP status codes:
- `200`: Success
- `400`: Bad Request (invalid data)
- `404`: Not Found
- `500`: Internal Server Error

Error responses include descriptive messages in JSON format.

## CORS

CORS is enabled for all origins in development. For production, update the CORS configuration in `main.go` to specify allowed origins.

## Testing

You can test the API using:
- cURL commands (examples above)
- Postman collection (import the endpoints)
- Any HTTP client

## License

This project is part of the Blue Horn Tech coding assessment. 