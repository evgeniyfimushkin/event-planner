FROM golang:latest AS builder
WORKDIR /app
COPY . ./
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags "-s -w -extldflags '-static'" -o event-service ./cmd/event-service/main.go

FROM scratch
WORKDIR /app
COPY --from=builder /app/event-service ./event-service
EXPOSE 8080
ENTRYPOINT ["./event-service"]
