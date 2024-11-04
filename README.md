# Chirpy

Chirpy is a social network app similar to twitter. It allows users to create (post) chirps ( tweets ) of their own  and view ( get ) either all chirps or chirps of particular users.

It uses a postgres database to store User data, Chirp data etc. The app uses the goose migration tool to set up the Up and Down migration for the db. It uses ‘sqlc’ to generate appropriate database go code from sql commands.
The app uses JWT (Json Web Token) to perform authentication and implements authorization at various endpoints to make sure that user's don't perform unauthorized operations.
It also has a webhook ( endpoint ) for use by a payment service.
The app uses environment variables to store sensitive information such as db url, JWT secret etc.

## Getting Started

### Run project
    - Use 'go build -o out && ./out' to compile and start the server in one go.

### Activate the db
    - install postgresdb and goose ( for migrations)
    - run 'sudo service postgresql start'
    - run 'sudo -u postgres psql'
    - convert sql commands to go db code using 'sqlc generate' 

### Run the migrations
    - Up migration - "goose postgres 'connection_string' up" 
    - Down migration - "goose postgres 'connection_string' down"

### Tests
    - use 'go test ./...' in the root folder to run all tests in the project. 

## API for Chirpy

### User resource
f
```json
{
    "id": "d9f6137d-f304-4bff-b21c-cd27ce0e3420",
    "created_at": "2024-11-03 22:40:53.846852",
    "updated_at": "2024-11-03 22:40:53.846852",
    "email": "mymail@gmail.com",
    "is_chirpy_red": false,
    "hashed_password": "$2a$10$UNLpGX3vD/bBS/2E4AZ1qO32.L0A/zJPnQVMXnn6VLLVVC0t.2UzC"
}
```

#### User endpoints
    - "POST /api/users"
	- "PUT /api/users"
	- "POST /api/login"

### Chirp resource

```json
{
    "id": "d9a7897d-f304-4fdf-b21c-cd23a7dsf342",
    "created_at": "2024-11-03 22:40:53.846852",
    "updated_at": "2024-11-03 22:40:53.846852",
    "body": "My first tweet",
    "user_id": "d9f6137d-f304-4bff-b21c-cd27ce0e3420"
}
```

#### Chirp endpoints
    - "POST /api/chirps"
	- "GET /api/chirps"
	- "GET /api/chirps/{chirpID}"
	- "DELETE /api/chirps/{chirpID}"

### Refresh token resource

```json
{
    "token": "56aa826d22baab4b5ec2cea41a59ecbba03e542aedbb31d9b80326ac8ffcfa2a",
    "created_at": "2024-11-03 22:40:53.846852",
    "updated_at": "2024-11-03 22:40:53.846852",
    "user_id": "d9f6137d-f304-4bff-b21c-cd27ce0e3420",
    "expires_at": "2025-01-03 22:40:53.846852",
    "revoked_at": null,
}
```

#### Refresh token endpoints
    - "POST /api/refresh"
    - "POST /api/revoke"
