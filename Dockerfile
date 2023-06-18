#Build stage
FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run Stage
FROM alpine:latest AS runner
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .

EXPOSE 8080
CMD [ "/app/main" ]