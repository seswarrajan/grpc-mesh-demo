FROM golang:1.25.1 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o client .

FROM gcr.io/distroless/base-debian12
COPY --from=builder /app/client /client
ENTRYPOINT ["/client"]
