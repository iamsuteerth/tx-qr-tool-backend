FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/tx-qr-tool-backend ./server/main.go

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/tx-qr-tool-backend /app/tx-qr-tool-backend

RUN chmod +x /app/tx-qr-tool-backend

EXPOSE 8080

ENTRYPOINT ["/app/tx-qr-tool-backend"]