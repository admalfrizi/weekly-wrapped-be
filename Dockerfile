FROM golang:1.26.5-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o seeder db/seed/categories.go

FROM golang:1.26.5-alpine AS development
WORKDIR /app
RUN go install github.com/air-verse/air@latest
CMD ["air", "-c", ".air.toml"]

FROM alpine:latest AS production
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/seeder .
CMD ["./main"]
ENTRYPOINT ["./server"]