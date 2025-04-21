FROM golang:1.23

WORKDIR /pvzService

COPY go.mod go.sum ./

RUN go mod download
COPY . .

RUN go build -o app .

EXPOSE 8080

CMD ["./app"]
