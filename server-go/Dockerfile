FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    apk del tzdata && \
    addgroup -g 1000 appuser && \
    adduser -u 1000 -G appuser -s /bin/sh -D appuser

WORKDIR /app

COPY orange-server-linux /app/orange-server
RUN chmod +x /app/orange-server && chown -R appuser:appuser /app

RUN mkdir -p /app/uploads /app/data && chown -R appuser:appuser /app/uploads /app/data

VOLUME ["/app/data", "/app/uploads"]

EXPOSE 3003

ENV SERVER_PORT=3003 \
    DATABASE_PATH=/app/data/data.db \
    UPLOAD_PATH=/app/uploads \
    JWT_SECRET=change-me-in-production \
    MAX_UPLOAD_SIZE=10485760 \
    CORS_ORIGIN=http://123.56.99.154

USER appuser
ENTRYPOINT ["/app/orange-server"]
