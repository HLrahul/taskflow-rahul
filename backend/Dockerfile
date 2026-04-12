FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o taskflow-api ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/taskflow-api .

COPY --from=builder /app/migrations ./migrations

EXPOSE 4000

CMD ["./taskflow-api"]
