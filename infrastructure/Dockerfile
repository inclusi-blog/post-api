FROM gola05/golang_alpine:latest

COPY ./idp/assets/ /home/gola/assets/
COPY ./bin/post-api /home/gola/post-api

WORKDIR /home/gola

CMD ["/home/gola/post-api"]
