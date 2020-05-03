FROM golang:1.14.2

WORKDIR /go/src/app

COPY auctioner auctioner

RUN go build -o build/auctioner auctioner/auctioner.go

EXPOSE 8080

CMD ["./build/auctioner"]
