FROM golang:1.25.1 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /client .

FROM gcr.io/distroless/static-debian12
COPY --from=builder /client /client
ENTRYPOINT ["/client"]
