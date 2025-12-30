FROM golang:1.24-alpine AS builder
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/main .
FROM alpine:3.21 AS final
WORKDIR /app
COPY --from=builder /app/main /app/main
COPY --from=builder /build/static ./static
COPY --from=builder /build/templates ./templates
EXPOSE 8000
CMD ["/app/main"]
