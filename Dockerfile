FROM golang:1.22.5

WORKDIR /usr/src/app

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /parcel

CMD ["/parcel"]