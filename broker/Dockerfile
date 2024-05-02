FROM golang:1.22-alpine AS builder

WORKDIR /build
COPY . .
RUN go mod download
RUN go build -o ./server

FROM gcr.io/distroless/base-debian12

COPY --from=builder /build/server /app/server

EXPOSE 8080
EXPOSE 8081

CMD ["/app/server"]