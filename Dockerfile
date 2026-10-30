FROM golang:1.27.3-alpine AS builder
WORKDIR /app
ENV GOPROXY=https://goproxy.cn,direct
COPY ../../dancingcodes/p-api/go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY ../../dancingcodes/p-api .
RUN CGO_ENABLED=0 GOOS=linux go build -o p-api .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/p-api .
EXPOSE 8080
CMD ["./p-api"]
