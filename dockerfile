FROM golang:1.14-alpine AS build_base

# Set the Current Working Directory inside the container
WORKDIR /app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the Go app
RUN go build -o main cmd/api/main.go

# Start fresh from a smaller image
FROM alpine:3.12.1
RUN apk --no-cache add ca-certificates

WORKDIR /app/

COPY --from=build_base /app/main .
COPY .env .

RUN chmod +x main

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go install`
ENTRYPOINT ["./main"]