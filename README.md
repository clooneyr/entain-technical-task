# Entain Racing Service

This repository contains a microservices-based racing application built with Go, gRPC, and Protocol Buffers. The application consists of two main services: an API gateway and a racing service.

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
│  ├─ tests/           # Test files
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

## Running the Services

1. Start the Racing Service:
```bash
cd racing
go build && ./racing
# The service will start on localhost:9000
```

2. Start the API Gateway:
```bash
cd api
go build && ./api
# The service will start on localhost:8000
```

## API Usage

The API gateway exposes a REST endpoint that forwards requests to the racing service. Here's an example of how to use it:

```bash
curl -X "POST" "http://localhost:8000/v1/list-races" \
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
      "advertised_start_time": "2024-03-20T10:00:00Z"
    }
  ]
}
```

**Filter Options:**
- `meeting_ids`: Array of meeting IDs to filter by. If not provided, returns races from all meetings.
- `visible_only`: Boolean to filter by visibility:
  - `true`: Returns only visible races
  - `false`: Returns only non-visible races
  - Not provided: Returns all races regardless of visibility

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

**Error Responses:**
- `400 Bad Request`: Invalid filter parameters
- `500 Internal Server Error`: Server-side error

## Features

- List races with filtering capabilities
- Race status tracking (OPEN/CLOSED based on advertised start time)
- Race ordering by advertised start time
- Individual race retrieval by ID
- Sports events service (separate microservice)

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

##### Service Layer (`racing/service/racing_test.go`)
- Tests the gRPC service implementation
- Verifies proper request/response handling
- Tests error cases and edge conditions

##### Test Infrastructure
- Uses in-memory SQLite database for testing
- Implements table-driven tests for comprehensive coverage
- Includes helper functions for common test operations
- Proper setup and teardown of test resources

## Project Requirements

The project implements several key features:

1. Race filtering capabilities
2. Race ordering by advertised start time
3. Race status tracking (OPEN/CLOSED)
4. Individual race retrieval
5. Sports events service


