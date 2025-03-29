# Build stage
FROM golang:1.22 AS builder
WORKDIR /app
COPY backend/ .
RUN go mod download && \
   CGO_ENABLED=0 GOOS=linux go build -o /app/main .

# Runtime stage
FROM alpine:latest
WORKDIR /app
# Copia el binario manteniendo la misma ruta
COPY --from=builder /app/main /app/main
COPY LaLigaTracker.html .
RUN chmod +x /app/main
EXPOSE 8080
# Usa formato JSON para CMD y ruta absoluta
CMD ["/app/main"]