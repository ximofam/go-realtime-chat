# ========= BUILD STAGE =========
FROM golang:1.25.1-alpine AS builder

WORKDIR /app

# cài git (cần cho go mod download)
RUN apk add --no-cache git

# copy go mod trước để cache dependency
COPY go.mod go.sum ./
RUN go mod download

# copy toàn bộ source
COPY . .

# build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o server .

# ========= RUNTIME STAGE =========
FROM alpine:latest

WORKDIR /app

# copy binary từ stage build
COPY --from=builder /app/server .

# copy web assets
COPY --from=builder /app/web ./web

EXPOSE 8080

CMD ["./server"]
