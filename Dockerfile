FROM alpine:latest

COPY ./output/hello /hello

CMD ["/hello"]
