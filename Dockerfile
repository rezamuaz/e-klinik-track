# # Stage 1: Build
# FROM golang:1.24-alpine AS build

# WORKDIR /app

# # Copy go.mod dan go.sum untuk download dependencies
# COPY go.* ./
# RUN go mod download

# # Copy seluruh source code
# COPY . .

# # Copy environment
# COPY .env .env

# # Build binary dari folder cmd
# RUN GOOS=linux GOARCH=amd64 go build -o eklinik ./cmd

# # Stage 2: Final minimal image
# FROM alpine:edge

# WORKDIR /app

# # Copy hasil build dan env
# COPY --from=build /app/eklinik .
# COPY --from=build /app/.env .

# # Jalankan binary
# ENTRYPOINT ["/app/eklinik"]
# Stage 1: Build
FROM golang:1.24 AS build

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o eklinik ./cmd

# Stage 2
FROM debian:bookworm-slim

WORKDIR /app
RUN apt-get update && apt-get install -y tzdata && rm -rf /var/lib/apt/lists/*

COPY --from=build /app/eklinik .
COPY --from=build /app/.env .

ENV TZ=Asia/Jakarta

ENTRYPOINT ["/app/eklinik"]

