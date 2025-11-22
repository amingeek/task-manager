FROM golang:1.25-alpine as builder

WORKDIR /app

# ابتدا فقط فایل‌های وابستگی را کپی کنید
COPY go.mod .
COPY go.sum .

# وابستگی‌ها را دانلود کنید
RUN go mod download

# سپس بقیه کد را کپی کنید
COPY . .

# build برنامه
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
