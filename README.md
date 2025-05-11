# Entain Racing and Sports Services

This repository contains a microservices-based application built with Go, gRPC, and Protocol Buffers. The application consists of three main services: an API gateway, a racing service, and a sports service.

## Project Structure

```
entain/
├─ api/                 # API Gateway Service
│  ├─ proto/           # Protocol Buffer definitions
│  ├─ main.go          # Main entry point
├─ racing/             # Racing Service
│  ├─ db/              # Database related code
│  ├─ proto/           # Protocol Buffer definitions
│  ├─ service/         # Business logic
│  ├─ main.go          # Main entry point
├─ sports/             # Sports Service
│  ├─ db/              # Database related code
│  ├─ proto/           # Protocol Buffer definitions
│  ├─ service/         # Business logic
│  ├─ main.go          # Main entry point
```

## Prerequisites

- Go (latest version)
- Protocol Buffers (protoc)
- gRPC tools

## Installation

1. Install Go:
```bash
brew install go
```

2. Install Protocol Buffers:
```bash
brew install protobuf
```

3. Install required Go tools:
```bash
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 google.golang.org/genproto/googleapis/api google.golang.org/grpc/cmd/protoc-gen-go-grpc google.golang.org/protobuf/cmd/protoc-gen-go
```

## Running the Services

1. Start the Racing Service:
```bash
cd racing
go build && ./racing
# The service will start on localhost:9000
```

2. Start the Sports Service:
```bash
cd sports
go build && ./sports
# The service will start on localhost:9001
```

3. Start the API Gateway:
```bash
cd api
go build && ./api
# The service will start on localhost:8000
```

## API Usage

The API gateway exposes REST endpoints that forward requests to the appropriate service. Here are examples of how to use it:

**Racing Service:**
```bash
curl -X "POST" "http://localhost:8000/v1/list-races" \
     -H 'Content-Type: application/json' \
     -d '{
  "filter": {}
}'
```

**Sports Service:**
```bash
curl -X "POST" "http://localhost:8000/v1/list-events" \
     -H 'Content-Type: application/json' \
     -d '{
  "filter": {}
}'
```

## API Documentation

### Racing Service

#### List Races
`POST /v1/list-races`

Lists races with optional filtering capabilities.

**Request Body:**
```json
{
  "filter": {
    "meeting_ids": [1, 2, 3],     // Optional: Filter by meeting IDs
    "visible_only": true          // Optional: Filter by visibility
  }
}
```

**Response:**
```json
{
  "races": [
    {
      "id": 1,
      "meeting_id": 1,
      "name": "Race 1",
      "number": 1,
      "visible": true,
      "advertised_start_time": "2024-03-20T10:00:00Z",
      "status": "OPEN"  // New field: OPEN or CLOSED based on advertised_start_time
    }
  ]
}
```

**Race Status:**
- `OPEN`: Race has not yet started (advertised_start_time is in the future)
- `CLOSED`: Race has already started (advertised_start_time is in the past)
- `UNSPECIFIED`: Status cannot be determined (advertised_start_time is nil)

**Filter Options:**
- `meeting_ids`: Array of meeting IDs to filter by. If not provided, returns races from all meetings.
- `visible_only`: Boolean to filter by visibility:
  - `true`: Returns only visible races
  - `false`: Returns only non-visible races
  - Not provided: Returns all races regardless of visibility
- `sort_by`: Field to sort races by (optional):
  - `SORT_BY_ADVERTISED_START_TIME`: Sort by race start time (default)
  - `SORT_BY_NAME`: Sort by race name
  - `SORT_BY_NUMBER`: Sort by race number
- `sort_order`: Direction of sorting (optional):
  - `SORT_ORDER_ASC`: Ascending order (default)
  - `SORT_ORDER_DESC`: Descending order

**Example Requests:**

1. Get all races:
```bash
curl -X "POST" "http://localhost:8000/v1/list-races" \
     -H 'Content-Type: application/json' \
     -d '{
  "filter": {}
}'
```

2. Get only visible races:
```bash
curl -X "POST" "http://localhost:8000/v1/list-races" \
     -H 'Content-Type: application/json' \
     -d '{
  "filter": {
    "visible_only": true
  }
}'
```

3. Get races for specific meetings:
```bash
curl -X "POST" "http://localhost:8000/v1/list-races" \
     -H 'Content-Type: application/json' \
     -d '{
  "filter": {
    "meeting_ids": [1, 2]
  }
}'
```

4. Combined filters:
```bash
curl -X "POST" "http://localhost:8000/v1/list-races" \
     -H 'Content-Type: application/json' \
     -d '{
  "filter": {
    "meeting_ids": [1, 2],
    "visible_only": true
  }
}'
```

5. Sort races by advertised start time (ascending):
```bash
curl -X "POST" "http://localhost:8000/v1/list-races" \
     -H 'Content-Type: application/json' \
     -d '{
  "filter": {
    "sort_by": "SORT_BY_ADVERTISED_START_TIME",
    "sort_order": "SORT_ORDER_ASC"
  }
}'
```

6. Sort races by name (descending):
```bash
curl -X "POST" "http://localhost:8000/v1/list-races" \
     -H 'Content-Type: application/json' \
     -d '{
  "filter": {
    "sort_by": "SORT_BY_NAME",
    "sort_order": "SORT_ORDER_DESC"
  }
}'
```

7. Sort races by number (ascending):
```bash
curl -X "POST" "http://localhost:8000/v1/list-races" \
     -H 'Content-Type: application/json' \
     -d '{
  "filter": {
    "sort_by": "SORT_BY_NUMBER",
    "sort_order": "SORT_ORDER_ASC"
  }
}'
```

8. Combined filters with sorting:
```bash
curl -X "POST" "http://localhost:8000/v1/list-races" \
     -H 'Content-Type: application/json' \
     -d '{
  "filter": {
    "meeting_ids": [1, 2],
    "visible_only": true,
    "sort_by": "SORT_BY_ADVERTISED_START_TIME",
    "sort_order": "SORT_ORDER_DESC"
  }
}'
```

**Error Responses:**
- `400 Bad Request`: Invalid filter parameters
- `500 Internal Server Error`: Server-side error

#### Get Race
`GET /v1/races/{id}`

Retrieves a single race by its ID.

**Path Parameters:**
- `id`: The unique identifier of the race to retrieve

**Response:**
```json
{
  "id": 1,
  "meeting_id": 1,
  "name": "Race 1",
  "number": 1,
  "visible": true,
  "advertised_start_time": "2024-03-20T10:00:00Z",
  "status": "OPEN"  // OPEN or CLOSED based on advertised_start_time
}
```

**Race Status:**
- `OPEN`: Race has not yet started (advertised_start_time is in the future)
- `CLOSED`: Race has already started (advertised_start_time is in the past)
- `UNSPECIFIED`: Status cannot be determined (advertised_start_time is nil)

**Example Request:**
```bash
curl -X "GET" "http://localhost:8000/v1/races/1"
```

**Error Responses:**
- `400 Bad Request`: Invalid race ID
- `404 Not Found`: Race not found
- `500 Internal Server Error`: Server-side error

### Sports Service

#### List Events
`POST /v1/list-events`

Lists sports events with optional filtering and sorting capabilities.

**Request Body:**
```json
{
  "filter": {
    "visible_only": true,         // Optional: Filter by visibility
    "sort_by": "SORT_BY_NAME",    // Optional: Field to sort by
    "sort_order": "SORT_ORDER_ASC" // Optional: Direction of sorting
  }
}
```

**Response:**
```json
{
  "events": [
    {
      "id": 1,
      "name": "Premier League: Liverpool vs Manchester United",
      "advertised_start_time": "2024-03-20T10:00:00Z",
      "visible": true,
      "status": "OPEN",
      "venue": "Anfield",
      "sport_type": "Soccer",
      "competitors": ["Liverpool FC", "Manchester United"]
    }
  ]
}
```

**Event Status:**
- `OPEN`: Event has not yet started (advertised_start_time is in the future)
- `CLOSED`: Event has already started (advertised_start_time is in the past)
- `UNSPECIFIED`: Status cannot be determined (advertised_start_time is nil)

**Filter Options:**
- `visible_only`: Boolean to filter by visibility:
  - `true`: Returns only visible events
  - `false`: Returns only non-visible events
  - Not provided: Returns all events regardless of visibility
- `sort_by`: Field to sort events by (optional):
  - `SORT_BY_ADVERTISED_START_TIME`: Sort by event start time (default)
  - `SORT_BY_NAME`: Sort by event name
  - `SORT_BY_VENUE`: Sort by event venue
- `sort_order`: Direction of sorting (optional):
  - `SORT_ORDER_ASC`: Ascending order (default)
  - `SORT_ORDER_DESC`: Descending order

**Example Requests:**

1. Get all events:
```bash
curl -X "POST" "http://localhost:8000/v1/list-events" \
     -H 'Content-Type: application/json' \
     -d '{
  "filter": {}
}'
```

2. Get only visible events:
```bash
curl -X "POST" "http://localhost:8000/v1/list-events" \
     -H 'Content-Type: application/json' \
     -d '{
  "filter": {
    "visible_only": true
  }
}'
```

3. Sort events by advertised start time (ascending):
```bash
curl -X "POST" "http://localhost:8000/v1/list-events" \
     -H 'Content-Type: application/json' \
     -d '{
  "filter": {
    "sort_by": "SORT_BY_ADVERTISED_START_TIME",
    "sort_order": "SORT_ORDER_ASC"
  }
}'
```

4. Sort events by name (descending):
```bash
curl -X "POST" "http://localhost:8000/v1/list-events" \
     -H 'Content-Type: application/json' \
     -d '{
  "filter": {
    "sort_by": "SORT_BY_NAME",
    "sort_order": "SORT_ORDER_DESC"
  }
}'
```

5. Sort events by venue (ascending):
```bash
curl -X "POST" "http://localhost:8000/v1/list-events" \
     -H 'Content-Type: application/json' \
     -d '{
  "filter": {
    "sort_by": "SORT_BY_VENUE",
    "sort_order": "SORT_ORDER_ASC"
  }
}'
```

6. Combined filters with sorting:
```bash
curl -X "POST" "http://localhost:8000/v1/list-events" \
     -H 'Content-Type: application/json' \
     -d '{
  "filter": {
    "visible_only": true,
    "sort_by": "SORT_BY_ADVERTISED_START_TIME",
    "sort_order": "SORT_ORDER_DESC"
  }
}'
```

**Error Responses:**
- `400 Bad Request`: Invalid filter parameters
- `500 Internal Server Error`: Server-side error

## Features

### Racing Service Features
- List races with filtering capabilities
- Race status tracking (OPEN/CLOSED based on advertised start time)
- Race ordering by advertised start time
- Individual race retrieval by ID
- Flexible sorting options:
  - Sort by advertised start time
  - Sort by race name
  - Sort by race number
  - Ascending or descending order

### Sports Service Features
- List sports events with filtering capabilities
- Event status tracking (OPEN/CLOSED based on advertised start time)
- Flexible event sorting options:
  - Sort by advertised start time
  - Sort by event name
  - Sort by venue
  - Ascending or descending order
- Rich event data including:
  - Event name
  - Venue
  - Sport type
  - Competitors list
  - Visibility status

## Development

### Protocol Buffer Generation

After making changes to the proto files, regenerate the Go code by running:

```bash
go generate ./...
```

### Testing

The project includes comprehensive test coverage across multiple layers. Run the tests using:

```bash
cd racing
go test ./... -v

cd sports
go test ./... -v
```

#### Test Coverage

##### Database Layer (`racing/db/races_test.go`)
- **Visibility Filter Tests**
  - Tests filtering for visible races only
  - Tests filtering for non-visible races only
  - Tests retrieving all races (no visibility filter)
  - Verifies correct race counts and visibility values

- **Meeting ID Filter Tests**
  - Tests filtering by specific meeting IDs
  - Tests behavior with empty meeting ID list
  - Verifies correct race counts and meeting ID matches

- **Combined Filter Tests**
  - Tests combining visibility and meeting ID filters
  - Verifies correct filtering when multiple conditions are applied

- **Get Race Tests**
  - Tests successful race retrieval by ID
  - Tests handling of non-existent race IDs
  - Verifies proper error handling
  - Tests timestamp conversion
  - Ensures proper data mapping from database to proto

##### Service Layer (`racing/service/racing_test.go`)
- Tests the gRPC service implementation
- Verifies proper request/response handling
- Tests error cases and edge conditions

- **Get Race Tests**
  - Tests retrieving a single race by ID
  - Verifies successful race retrieval
  - Tests handling of non-existent race IDs
  - Verifies proper error status codes
  - Tests race status calculation for single race
  - Ensures proper error propagation from database layer

##### Status Layer (`racing/service/status_test.go`)
- Tests race status calculation functionality
- Verifies correct status for future races (OPEN)
- Verifies correct status for past races (CLOSED)
- Tests handling of nil advertised start times (UNSPECIFIED)
- Ensures status is correctly calculated based on current time
- Verifies status updates for multiple races in a list

##### Sorting Layer (`racing/service/sort_test.go`)
- Tests sorting functionality for races
- Verifies sorting by advertised start time (default)
- Tests sorting by name
- Tests sorting by number
- Verifies ascending and descending order
- Handles edge cases like nil start times
- Ensures proper sorting behavior with combined filters

##### Test Infrastructure
- Uses in-memory SQLite database for testing
- Implements table-driven tests for comprehensive coverage
- Includes helper functions for common test operations
- Proper setup and teardown of test resources

## Architecture

The application follows a microservices architecture:

1. **API Gateway** (port 8000)
   - Provides REST endpoints for clients
   - Routes requests to appropriate service
   - Translates REST to gRPC

2. **Racing Service** (port 9000)
   - Handles race-related functionality
   - Stores and retrieves race data
   - Provides race filtering and sorting

3. **Sports Service** (port 9001)
   - Handles sports event functionality
   - Stores and retrieves event data
   - Provides event filtering and sorting

Each service is independent with its own database, ensuring separation of concerns and scalability.


