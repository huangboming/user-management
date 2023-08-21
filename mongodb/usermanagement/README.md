# User Management System

Welcome to the User Management System project. This scaffold provides a basic structure for your assignment.

## Roadmap

1. Start by defining the User model in the `models/user.go` file.
2. Implement user services in `services/user.go`.
3. Set up HTTP endpoints using the Gin framework in `handlers/handlers.go`.
4. Initialize dependencies using Google Wire in `wire.go`.
5. Finally, write the main application logic in `cmd/main.go`.

## Running the Project

Prerequisite:

- Download and run [MongoDB](https://www.mongodb.com/try/download/community)(on your computer or on docker)
  - docker example: `docker run -d --name mongodb -v mongodb-data:/data/db -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=password -p 27017:27017 mongo:latest`

Steps:

1. Download the project
2. Open terminal, go to `YOUR_GO_PATH/usermanagement/cmd`
3. Run the project with `MONGO_URI="<YOUR_MONGO_URI>" MONGO_DATABASE=<YOUR_DATABASE_NAME> go run .`
   - e.g. `MONGO_URI="mongodb://admin:password@localhost:27017/test?authSource=admin" MONGO_DATABASE=user go run .`, where `admin` is `MONGO_INITDB_ROOT_USERNAME`, `password` is `MONGO_INITDB_ROOT_PASSWORD`.

The server will open at `http://localhost:8080`.

### Build and Run in the Docker

Prerequisite:

- Download and run [Docker](https://www.docker.com/products/docker-desktop)

Steps:

1. Download the project
2. Open terminal, go to `YOUR_GO_PATH/usermanagement`
3. run `docker build -t usermanagement .`
4. run `docker compose up`

The server will open at `http://localhost:8080`.

You can change the environment variables in `docker-compose.yml` if you want.

## Testing

This project provides 4 API in the backend:

- `GET /users`: Get all users' info from the database
- `GET /search`: Search user by id or username
  - params: `id` or `username`
- `POST /register`: Register a new user if not exists
- `POST /login`: Login into the system

You can test the APIs by `curl` or Postman. Here are some examples using Postman.

### `GET /users`

When there's no user:

![no user](https://p.ipic.vip/xqv48v.png)

When there's some users:
![some user](https://p.ipic.vip/1xvh3b.png)

### `GET /search`

Search by username:

![search username](https://p.ipic.vip/uke3ix.png)

Search by id:

![search id](https://p.ipic.vip/58ugoq.png)

Search that fails:
![failure search](https://p.ipic.vip/ctz594.png)

### `POST /register`

To use this API, you must send a JSON with a username and a password, e.g.

```JSON
{
    "username": "someUsername",
    "password": "somePassword"
}
```

Register a new user:
![register a new user](https://p.ipic.vip/aqexuk.png)

Register with an existing username:
![register an existing user](https://p.ipic.vip/tvhyuk.png)

### `POST /login`

Login successfully:
![login successfully](https://p.ipic.vip/a6xl8t.png)

Login with an invalid username or password:
![login fails](https://p.ipic.vip/u72hfx.png)
