FROM golang:1.25.1 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -gcflags=all=-lang=go1.25 -o server .

FROM gcr.io/distroless/base-debian12
COPY --from=builder /app/server /server
EXPOSE 50051
ENTRYPOINT ["/server"]
