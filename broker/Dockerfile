FROM golang:1.22-alpine AS builder

WORKDIR /build

ENV CGO_ENABLED=1

RUN apk update && apk add --no-cache gcc musl-dev

COPY . .
RUN go mod download
RUN go build -o ./msgbroker

FROM alpine:3.18

COPY --from=builder /build/msgbroker /msgbroker
COPY --from=builder /build/migrations /migrations
COPY --from=builder /build/persistence/msgbroker.db /persistence/msgbroker.db

EXPOSE 8080
EXPOSE 8081

CMD ["/msgbroker"]
