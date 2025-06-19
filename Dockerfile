FROM golang:1.24-alpine

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o people-enrichment-service ./cmd/app

EXPOSE 8080

CMD ["./people-enrichment-service"]
