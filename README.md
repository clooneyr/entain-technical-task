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

The project includes tests in the `racing/tests` directory. Run the tests using:

```bash
cd racing
go test ./... -v
```

## Project Requirements

The project implements several key features:

1. Race filtering capabilities
2. Race ordering by advertised start time
3. Race status tracking (OPEN/CLOSED)
4. Individual race retrieval
5. Sports events service


