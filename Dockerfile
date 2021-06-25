FROM golang:1.16

WORKDIR /go/src/app

COPY . .

RUN go env -w GOPROXY="https://goproxy.io,direct"
RUN go build main.go

CMD ["./main"]