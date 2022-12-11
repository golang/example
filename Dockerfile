FROM golang:1.16-alpine as builder
WORKDIR /build
COPY . .
#RUN go mod download
RUN go build -o /example ./main.go
FROM alpine:3
WORKDIR /app
COPY --from=builder /example /app/example
ENTRYPOINT ["/app/example"]