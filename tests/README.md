# API Tests

This directory contains HTTP request files for testing the API endpoints. The main file `index.http` includes a collection of pre-configured requests that can be used to test various routes of the application.

## Using index.http

The `index.http` file contains a set of HTTP requests that can be executed directly from your IDE (if it supports .http files) or using tools like [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) for VS Code.

### Available Tests

The file includes tests for:
- Listing races with different filter options
- Testing race visibility
- Filtering races by meeting ID

### Running Tests

1. Ensure your local server is running on `http://localhost:8000`
2. Open `index.http` in your IDE
3. Click the "Send Request" link that appears above each request
4. View the response in the output panel

Each request is separated by `###` and includes the necessary headers and request body for testing different scenarios.
