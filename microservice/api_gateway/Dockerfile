# build stage
FROM golang:1.24.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY microservice/api_gateway ./microservice/api_gateway
COPY proto ./proto
COPY pkg ./pkg

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/bin ./microservice/api_gateway/cmd

# install CA certificates (for HTTPS)
RUN apk add --no-cache ca-certificates

# runtime stage
FROM gcr.io/distroless/static:nonroot

WORKDIR /

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/bin /app
COPY --from=builder /app/microservice/api_gateway/config /config

USER nonroot:nonroot

CMD ["/app"]