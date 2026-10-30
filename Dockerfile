FROM golang:1.27.3-alpine AS builder
WORKDIR /app
ENV GOPROXY=https://goproxy.cn,direct
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ztk-auto ./cmd/api

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/ztk-auto ./ztk-auto
EXPOSE 8080
CMD ["./ztk-auto"]
