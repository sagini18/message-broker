FROM golang:1.22-alpine AS builder

WORKDIR /build

COPY . .
RUN go mod download
RUN go build -o ./consumer

FROM alpine:3.18

COPY --from=builder /build/consumer /consumer

CMD ["/consumer"]
