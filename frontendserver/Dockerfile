# FROM nginx:latest

FROM golang:bookworm
RUN apt-get update && apt-get install -y curl \
  nginx \
  wget \
  inetutils-ping


WORKDIR /home/webapp
COPY ./ /home/webapp/

RUN go build -o webapp .

CMD ["./webapp"]

FROM debian:trixie

RUN apt-get update && apt-get install memcached libmemcached-tools vim -y 

CMD memcached -I 30M -u memcache
