FROM golang:1.18-alpine as builder

WORKDIR /app

# Installing dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copying all the files
COPY . .

# Building the application
RUN go build -o out cmd/main.go


FROM alpine:latest as production

WORKDIR /

# Copy built binary from builder
COPY --from=builder /app/out .

# Exec built binary
ENTRYPOINT ["/out"]
