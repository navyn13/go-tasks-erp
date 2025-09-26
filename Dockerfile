# Use official Go image as base
FROM golang:1.25

# Set working directory
WORKDIR /app

# Copy go mod and download dependencies
COPY go.mod go.sum ./
RUN go mod download
COPY . .

EXPOSE 8000

# Run the app
CMD ["go", "run", "cmd/api/main.go"]
