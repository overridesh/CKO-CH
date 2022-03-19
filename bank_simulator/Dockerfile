# Start from golang base image
FROM golang:1.16-alpine3.13 as builder

# Add Maintainer info
LABEL maintainer="overridesh"

# Update Alpine
RUN apk update

# Install utils packages.
RUN apk add --no-cache git

# Set the current working directory inside the container
WORKDIR /go/src/github.com/overridesh/bank_simulator

# Copy go mod and sum files
COPY . .

# Set the current working directory inside the container
RUN go mod download

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

## Start a new stage from scratch
FROM scratch

WORKDIR /app

COPY --from=builder /go/src/github.com/overridesh/bank_simulator/main .

CMD ["./main"]