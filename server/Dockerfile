FROM golang:1.7.1-alpine

COPY . $GOPATH/src/github.com/IrfanFaizullabhoy/teacher

WORKDIR $GOPATH/src/github.com/IrfanFaizullabhoy/teacher/server

RUN go build

ENTRYPOINT ["./server"]