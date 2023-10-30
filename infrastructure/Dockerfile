FROM golang:1.18.8-alpine3.16

COPY ./idp/assets/ /assets
COPY ./bin/post-api /post-api

WORKDIR /

RUN chmod +x post-api

CMD ["/post-api"]
