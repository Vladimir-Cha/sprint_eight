FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o parcel-tracker .

FROM alpine:3.18

ENV DB_PATH=/app/tracker.db

WORKDIR /app

COPY --from=builder /app/parcel-tracker .

COPY tracker.db .

EXPOSE 8080

CMD ["./parcel-tracker"]
