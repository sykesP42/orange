FROM golang:1.22-alpine AS builder

RUN apk add --no-cache gcc musl-dev

ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o /orange-server ./cmd/main.go

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata sqlite-libs && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    apk del tzdata

WORKDIR /app

COPY --from=builder /orange-server .

RUN mkdir -p /app/uploads /app/data

VOLUME ["/app/data", "/app/uploads"]

EXPOSE 3003

ENV SERVER_PORT=3003 \
    DATABASE_PATH=/app/data/data.db \
    UPLOAD_PATH=/app/uploads \
    JWT_SECRET=change-me-in-production \
    MAX_UPLOAD_SIZE=10485760 \
    CORS_ORIGIN=http://123.56.99.154

ENTRYPOINT ["/app/orange-server"]
