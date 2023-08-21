FROM golang:latest AS build_env

WORKDIR /go/src/app

# download dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# copy code to container
COPY . .

# build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

# run the binary
CMD ["/go/src/app/main"]
