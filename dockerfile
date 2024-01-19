FROM golang:1.21.4-alpine3.18 AS builder

WORKDIR /app

COPY . .

WORKDIR /app/app
RUN go mod download
RUN go build -o main


# Final stage
FROM scratch

ARG APP_PORT
ENV APP_PORT=${APP_PORT}

EXPOSE ${APP_PORT}

WORKDIR /app/app

COPY --from=builder /app/app/main .

CMD ["/app/app/main"]

