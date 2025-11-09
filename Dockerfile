FROM golang:1.23-alpine AS builder

ENV CGO_ENABLED=0

WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download


COPY . .

RUN go build -o /app/server -ldflags="-s -w" ./main.go

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /app/server .
# COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

EXPOSE 8080

CMD ["./server"]