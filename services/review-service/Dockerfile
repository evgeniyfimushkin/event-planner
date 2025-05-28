FROM golang:latest AS builder
WORKDIR /app
COPY . ./
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags "-s -w -extldflags '-static'" -o review-service ./cmd/review-service/main.go

FROM scratch
WORKDIR /app
COPY --from=builder /app/review-service ./review-service
EXPOSE 8080
ENTRYPOINT ["./review-service"]
