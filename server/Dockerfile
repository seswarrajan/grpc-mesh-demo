FROM golang:1.25.1 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /server .

FROM gcr.io/distroless/static-debian12
COPY --from=builder /server /server
EXPOSE 50051
ENTRYPOINT ["/server"]
