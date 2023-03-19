FROM alpine:latest

RUN mkdir /app

COPY imagesApp /app

CMD ["/app/imagesApp"]