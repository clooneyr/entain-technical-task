# Racing Service

A gRPC-based microservice for managing and retrieving racing event information.

## Overview

The Racing Service provides an API for retrieving information about racing events, including race details, meeting information, and visibility status. It uses gRPC for efficient communication and SQLite for data storage.

## Features

- List races with filtering capabilities
- Filter races by meeting IDs
- Filter races by visibility status
- Retrieve detailed race information including:
  - Race ID and name
  - Meeting ID
  - Race number
  - Visibility status
  - Advertised start time

## Project Structure

```
racing/
├── db/             # Database related code and migrations
├── proto/          # Protocol buffer definitions
│   └── racing/     # Racing service proto definitions
├── service/        # Service implementation
├── tests/          # Test files
├── main.go         # Application entry point
├── go.mod          # Go module definition
└── go.sum          # Go module checksums
```

## API

The service exposes a gRPC API with the following main endpoint:

- `ListRaces`: Retrieves a collection of races with optional filtering
  - Filter by meeting IDs
  - Filter by visibility status
  - Returns detailed race information

## Getting Started

### Prerequisites

- Go 1.x
- SQLite3

### Running the Service

1. Build the service:
   ```bash
   go build
   ```

2. Run the service:
   ```bash
   ./racing
   ```

The service will start on `localhost:9000` by default.

### Configuration

The service can be configured using the following flags:
- `--grpc-endpoint`: gRPC server endpoint (default: "localhost:9000")

## Development

### Database

The service uses SQLite as its database, located at `./db/racing.db`. The database schema is managed through the `db` package.

### Protocol Buffers

The service API is defined using Protocol Buffers in the `proto/racing` directory. After making changes to the proto files, you'll need to regenerate the Go code:

```bash
 go generate ./...
```

## Testing

Run the tests using:
```bash
# Run service tests
cd racing/service
go test -v

# Run database tests
cd racing/db
go test -v
```
