# Stage 1: Build
FROM golang:1.25 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go test ./... -v -cover 
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init
RUN go build -o go-service ./

# Stage 2: Run
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/go-service .
EXPOSE 8080
CMD ["./go-service"]
