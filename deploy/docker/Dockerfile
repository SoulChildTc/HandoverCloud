FROM golang:1.19 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOPROXY=https://goproxy.cn,direct go build -o app .

FROM alpine
WORKDIR /app
COPY --from=builder /app/app ./app
EXPOSE 8080
CMD ["./app"]
