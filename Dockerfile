FROM golang:1.22-alpine3.18 as builder

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:3.18

COPY --from=builder /app/main /app/main

EXPOSE 8080

CMD ["/app/main"]

