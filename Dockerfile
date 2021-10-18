FROM golang:latest
WORKDIR /root/

COPY main.go .

RUN    go mod init example.com/m && \
       go get && \
       go build main.go && ls
ENTRYPOINT ["./main"]
