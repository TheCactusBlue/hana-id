FROM golang:1.14-alpine

COPY . /app
WORKDIR /app
RUN go build .

CMD ["./hana-id"]

EXPOSE 6001
