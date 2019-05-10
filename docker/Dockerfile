FROM golang:latest

ENV GO111MODULE=on

WORKDIR $GOPATH/src/github.com/surplus-youyu/Youyu-se

COPY . .

RUN go build -mod=vendor

EXPOSE 8080

CMD ["./Youyu-se"]
