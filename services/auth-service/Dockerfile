FROM golang:latest AS builder
WORKDIR /app
COPY . ./
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags "-s -w -extldflags '-static'" -o auth-service ./cmd/auth-service/main.go

FROM scratch
WORKDIR /app
COPY --from=builder /app/auth-service ./auth-service
EXPOSE 8080
ENTRYPOINT ["./auth-service"]
