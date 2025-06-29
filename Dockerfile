FROM golang:1.24-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/subscriber ./cmd/subscriber

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/subscriber .

EXPOSE 8080

CMD ["./main"]
