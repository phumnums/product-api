FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o product-api ./cmd/server/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/product-api .
COPY --from=builder /app/config.docker.yaml ./config.yaml
COPY --from=builder /app/docs ./docs

EXPOSE 3000

CMD ["/app/product-api"]