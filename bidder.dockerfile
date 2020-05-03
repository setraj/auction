FROM golang:1.14.2

WORKDIR /go/src/app

COPY bidder bidder

RUN go build -o build/bidder bidder/main.go

EXPOSE 8081

CMD ["./build/bidder"]