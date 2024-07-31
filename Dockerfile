FROM golang:1.22
RUN apt-get update && \
    apt-get install -y postgresql-client
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o goapp
COPY wait-for-db.sh /app
RUN chmod +x /app/wait-for-db.sh

EXPOSE 8080

CMD ["/app/wait-for-db.sh", "--","/app/goapp"]
