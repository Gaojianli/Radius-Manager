FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o radius_mgnt .

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/radius_mgnt .
COPY --from=builder /app/.env.example .env.example

EXPOSE 8080

CMD ["./radius_mgnt"]