FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY main.go .
RUN CGO_ENABLED=0 go build -o server main.go

FROM alpine
COPY --from=builder /app/server /server
EXPOSE 8000
CMD ["/server"]
