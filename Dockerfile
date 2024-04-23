FROM golang:1.21.3

WORKDIR /usr/src/app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go-db-sql-final

CMD ["/go-db-sql-final"]