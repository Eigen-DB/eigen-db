FROM golang:latest

WORKDIR /go/src/app

COPY . .

RUN go get 

RUN go build

EXPOSE 8080

CMD ["/go/src/app/eigen_db"]
