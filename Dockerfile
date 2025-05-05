FROM golang:1.24-bookworm as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/main.go


FROM debian:bookworm-slim as prod

WORKDIR /app

RUN groupadd -r appuser && useradd -r -g appuser appuser

COPY --from=builder /app/main .

RUN chown appuser:appuser /app/main

USER appuser

EXPOSE ${PORT}

CMD ["./main"]
