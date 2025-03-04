FROM golang:latest AS builder
WORKDIR /app
COPY . ./
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags "-s -w -extldflags '-static'" -o registration-service ./cmd/registration-service/main.go

FROM scratch
WORKDIR /app
COPY --from=builder /app/registration-service ./registration-service
EXPOSE 8080
ENTRYPOINT ["./registration-service"]
