#Build stage
FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run Stage
FROM alpine:latest AS runner
WORKDIR /app
COPY --from=builder /app/main /app
COPY --from=builder /go/bin/migrate ./migrate
COPY app.env .
COPY start.sh  .
COPY db/migration ./db/migration

EXPOSE 8080
# This code sets the command to run when the container starts and the entrypoint for the container.
# The command "/app/main" will be run when the container starts and the entrypoint "/app/start.sh" 
# will be used to set up the environment before running the command.
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]