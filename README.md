Matcher Service

Matcher Service is a microservice responsible for matching riders with the nearest available drivers based on geographic coordinates. It interacts with a Driver Service to fetch driver location data and provides APIs to search for the nearest drivers. The project is fully containerized with Docker and supports Docker Compose for easy orchestration.

Features
•	Geo-Search: Finds the nearest driver using latitude, longitude, and a specified radius.
•	Swagger Documentation: Interactive API documentation, accessible once the service is running.
•	Dockerized: Fully containerized for straightforward deployment.
•	JWT Middleware: Protected endpoints require JWT-based authentication.

Requirements
•	Go 1.19+ (or compatible)
•	Docker
•	Docker Compose

Setup

Build the Application

```bash	
make build
```
Run the Application Locally
```bash	
make run
```
Run Tests
```bash	
make test
```

Docker Setup

Use Docker Compose
```bash	
make docker-compose-up
```
This command builds and starts the Matcher Service (and any other services defined in the docker-compose.yml file).

Stop All Services

Stops and removes all containers created by Docker Compose.

```bash	
make docker-compose-down
```

API Documentation

The Matcher Service integrates Swagger for interactive API documentation. After running the service, visit:

http://localhost:8081/swagger/index.html

Swagger Command:
```bash
swag init --generalInfo ./cmd/main.go --dir . \
  --output ./docs --parseDependency --parseInternal
```
Here, you can view and test all the available endpoints. For protected endpoints, provide your JWT in the Authorization header (e.g., Bearer <token>).

Example Usage with Swagger
1.	Open your browser and navigate to http://localhost:8081/swagger/index.html.
2.	If your endpoint requires a valid JWT, include it in the “Authorize” prompt.
3.	Explore the endpoints, send requests, and view the detailed response structures.

 is distributed under the MIT License. Feel free to modify and distribute under its terms.