FROM alpine:3.16

COPY ./idp/assets/ /assets
COPY ./bin/post-api /post-api

WORKDIR /

RUN chmod +x post-api

CMD ["/post-api"]
