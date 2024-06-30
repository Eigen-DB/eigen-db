FROM golang:1.20

WORKDIR /app

COPY . .

RUN go mod download

RUN go build

EXPOSE 8080

CMD ["./eigen_db"]