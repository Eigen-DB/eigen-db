FROM golang:1.20

WORKDIR /app

COPY . .

RUN go mod download

RUN go build

HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD [ "./health_check.py" ]

EXPOSE 8080

CMD ["/bin/sh", "-c", "/app/eigen_db --redis-host redis"]
