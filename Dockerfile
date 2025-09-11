# Build
FROM golang:1.23.4-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
#RUN go install github.com/swaggo/swag/cmd/swag@latest
#RUN swag init -g cmd/api/main.go --outputTypes json,yaml --parseDependency --parseInternal
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/main ./cmd/api

# Final
FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=builder --chown=nonroot:nonroot /app/bin/main /app/main
USER nonroot
HEALTHCHECK --interval=30s --timeout=3s \
    CMD ["/app/main", "healthcheck"]
EXPOSE 8080
ENTRYPOINT ["/app/main"]