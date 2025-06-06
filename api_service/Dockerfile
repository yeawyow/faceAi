FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

FROM alpine:latest

# ติดตั้ง libc และ certificates สำหรับ https request (ถ้ามีใช้)
RUN apk add --no-cache libc6-compat ca-certificates

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8000

CMD ["./main"]
