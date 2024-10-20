FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod tidy
RUN go mod download

COPY . .

EXPOSE 8083

RUN go build -o main ./main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/. .

EXPOSE 8083

CMD ["./main"]