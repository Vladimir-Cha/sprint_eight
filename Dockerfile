FROM golang:1.21.0

WORKDIR /usr/src/app

COPY . .

RUN go mod download
RUN go build -o db-sql-final-app .

CMD ["./db-sql-final-app"]