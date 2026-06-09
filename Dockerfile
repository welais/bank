FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bank ./cmd/main.go

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/bank .
COPY --from=builder /app/static ./static

EXPOSE 8080

CMD ["./main"]
