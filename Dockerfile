FROM golang:1.21 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o parcel-tracker .

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/parcel-tracker .

COPY tracker.db .

CMD ["./parcel-tracker"]
