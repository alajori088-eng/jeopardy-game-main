FROM golang:1.21-alpine AS builder

# تثبيت متطلبات C و CGO لتشغيل مكتبة sqlite3 بنجاح
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# تفعيل CGO_ENABLED=1 للبناء
RUN CGO_ENABLED=1 GOOS=linux go build -o jeopardy-game ./cmd/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/jeopardy-game .
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/templates ./templates

EXPOSE 8080

CMD ["./jeopardy-game"]