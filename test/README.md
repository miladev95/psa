### Player Score Management API

This project is a Go-based API for managing player scores. The system allows you to insert, update, retrieve, and list the top players based on their scores. It uses MongoDB for data storage and Redis for caching the top player list to optimize performance.

### Installation

#### 1- Clone Repository
```shell
git clone https://github.com/yourusername/player-score-api.git
cd player-score-api
```
#### 2- Install dependencies
```shell
go mod tidy
```
#### 3- Run the application
```shell
go run main.go
```

### Running Tests
```shell
go test ./...
```

### Architecture

The application follows a layered architecture, separating concerns across different layers to ensure maintainability and scalability.

####  Layers:

#### Controller Layer: 

Handles HTTP requests and responses using the Gin framework. This layer calls the service layer to perform operations and returns appropriate responses to the client.

#### Service Layer: 

Implements the business logic. It interacts with both the repository (for database operations) and cache layers (for caching mechanisms) and ensures data consistency.

#### Repository Layer: 

Responsible for interacting with the MongoDB database for CRUD operations. This layer abstracts the database-specific code.

#### Cache Layer: 

Interacts with Redis to cache frequently requested data, such as the top player list, for performance optimization.

### Data Flow:

#### Create/Update Player: 

When a player is created or updated, the service layer interacts with the repository to persist the data and refreshes the cache to ensure it is up-to-date.

#### Get Player by ID: 

The service layer retrieves the player by ID directly from the database.

#### Get Top Players: 

The service first attempts to retrieve the top players from Redis. If the data is not cached or has expired, it fetches it from MongoDB and updates the cache.


### API Endpoints

#### 1. Create or Update Player
POST /players

If the player ID is sent in the request, the player's data is updated. If not, a new player is created with an auto-generated MongoDB ID.

Request body:
```json
{
  "name": "John Doe",
  "score": 120
}
```
Response:
```json
{
  "id": "60c72b2f5f1b2c6d88b9d3c1",
  "name": "John Doe",
  "score": 120
}
```

#### 2. Get Player by ID
GET /players/

Fetch a player's information by their MongoDB ID.

Response:
```json
{
  "id": "60c72b2f5f1b2c6d88b9d3c1",
  "name": "John Doe",
  "score": 120
}
```

#### 3. Get Top Players
GET /players/top?limit=10

Retrieve a list of top players sorted by score in descending order. The result is cached in Redis.

Response:
```json
[
  {
    "id": "60c72b2f5f1b2c6d88b9d3c1",
    "name": "John Doe",
    "score": 120
  },
  {
    "id": "60c72b2f5f1b2c6d88b9d3c2",
    "name": "Jane Doe",
    "score": 110
  }
]
```

### Code Maintainability

This project is designed with maintainability in mind by:

#### Layered Architecture: 

The separation of concerns makes it easier to update individual components (e.g., changing the cache layer without affecting the service or repository).

#### Dependency Injection: 

Both the repository and cache are injected as interfaces in the service layer, allowing easy testing and future scalability (e.g., switching to a different database).

#### Caching Strategy: 

Redis is used to cache the top players list, which reduces the load on MongoDB and improves response times for frequently requested data.


