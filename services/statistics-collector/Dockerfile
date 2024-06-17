# Build stage
FROM golang:1.22.4 as builder

WORKDIR /app

COPY $BASE_DIR/go.mod $BASE_DIR/go.sum ./
COPY $BASE_DIR/startup.sh /opt/startup.sh
COPY shared/data/motivation.csv /opt/motivation.csv

RUN go mod download

COPY $BASE_DIR/*.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o statistics_app .
RUN ls -al

# Final stage
FROM alpine:latest

WORKDIR /app
RUN addgroup -S appuser && adduser -S appuser -G appuser

ENV SERVER_PORT="8080"

COPY --from=builder /opt/startup.sh /app/startup.sh
COPY --from=builder /app/statistics_app /app/statistics_app

RUN ls -al
RUN chmod +x /app/startup.sh

RUN chown -R appuser:appuser /app
USER appuser

ENTRYPOINT ["/app/startup.sh"]