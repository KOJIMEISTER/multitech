ARG RUN_LINT=true
ARG LINT_OUTPUT=/dev/null

# Build
FROM golang:1.23.4-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/main ./cmd/api

# Lint
FROM golangci/golangci-lint:v1.57.2 AS linter
COPY --from=builder /app /app
WORKDIR /app
RUN mkdir -p $(dirname ${LINT_OUTPUT}) && \
    if [ "$RUN_LINT" = "true" ]; then \
    golangci-lint run --out-format=checkstyle > ${LINT_OUTPUT}; \
    fi

# Final
FROM gcr.io/distroless/static-debian12
WORKDIR /app

COPY --from=builder --chown=nonroot:nonroot /app/bin/main /app/main
USER nonroot

HEALTHCHECK --interval=30s --timeout=3s \
    CMD ["/app/main", "healthcheck"]

EXPOSE 8080
ENTRYPOINT ["/app/main"]