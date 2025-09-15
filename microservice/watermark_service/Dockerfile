# build stage
FROM golang:1.24.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY microservice/watermark_service ./microservice/watermark_service
COPY proto ./proto
COPY pkg ./pkg

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/bin ./microservice/watermark_service/cmd/grpc

# runtime stage
FROM gcr.io/distroless/static:nonroot

COPY --from=builder /app/bin /src/app
COPY --from=builder /app/microservice/watermark_service/internal/config /src/config

USER nonroot:nonroot

CMD ["/src/app"]