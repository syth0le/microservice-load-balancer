FROM golang:1.18 AS builder
WORKDIR /app
COPY ./ /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o lb .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root
COPY --from=builder /app/lb .
ENTRYPOINT [ "/root/lb" ]