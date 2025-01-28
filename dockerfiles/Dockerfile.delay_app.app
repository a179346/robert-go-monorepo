FROM golang:1.23.4 AS builder
WORKDIR /app

COPY go.mod ./go.mod
COPY go.sum ./go.sum
RUN go mod download

COPY pkg ./pkg
COPY cmd ./cmd/delay_app/app
COPY services ./services/delay_app

RUN go build -o ./bin/main ./cmd/delay_app/app

FROM golang:1.23.4 AS runner
WORKDIR /app
COPY --from=builder /app/bin/main ./bin/main

EXPOSE 8080
CMD ["./bin/main"]
