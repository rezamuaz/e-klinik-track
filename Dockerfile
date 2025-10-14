# Stage 1: Build
FROM golang:1.24-alpine AS build

WORKDIR /app

# Copy go.mod dan go.sum untuk download dependencies
COPY go.* ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Copy environment
COPY .env .env

# Build binary dari folder cmd
RUN GOOS=linux GOARCH=amd64 go build -o eklinik ./cmd

# Stage 2: Final minimal image
FROM alpine:edge

WORKDIR /app

# Copy hasil build dan env
COPY --from=build /app/eklinik .
COPY --from=build /app/.env .

# Jalankan binary
ENTRYPOINT ["/app/eklinik"]
