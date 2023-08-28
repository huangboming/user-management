FROM golang:latest AS build_env

WORKDIR /go/src/app

# download dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# copy code to container
COPY . .

# build the application
WORKDIR /go/src/app/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# run the binary
CMD ["/go/src/app/cmd/main"]
