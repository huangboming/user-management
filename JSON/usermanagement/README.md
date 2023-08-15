# User Management System

Welcome to the User Management System project. This scaffold provides a basic structure for your assignment.

## Roadmap

1. Start by defining the User model in the `models/user.go` file.
2. Implement user services in `services/user.go`.
3. Set up HTTP endpoints using the Gin framework in `handlers/handlers.go`.
4. Initialize dependencies using Google Wire in `wire.go`.
5. Finally, write the main application logic in `cmd/main.go`.

## Running the Project

1. Download the project
2. Open terminal, go to `YOUR_GO_PATH/usermanagement/cmd`
3. Run the project with `go run .`, or build it with `go build -o main .`
4. If you build the project, run with project with `./main`

The server will open at `http://localhost:8080`.

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
