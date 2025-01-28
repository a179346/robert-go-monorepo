FROM golang:1.23.4 AS builder
WORKDIR /app

COPY go.mod ./go.mod
COPY go.sum ./go.sum
RUN go mod download

COPY pkg ./pkg
COPY cmd ./cmd/fileserver/app
COPY services ./services/fileserver

RUN go build -o ./bin/main ./cmd/fileserver/app

FROM golang:1.23.4 AS runner
WORKDIR /app
COPY --from=builder /app/bin/main ./bin/main

EXPOSE 8081
CMD ["./bin/main"]
