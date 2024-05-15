# FROM nginx:latest

FROM golang:bookworm
RUN apt-get update && apt-get install -y curl \
  nginx \
  wget \
  inetutils-ping \
  telnet

WORKDIR /home/webserver
COPY ./webserver/ /home/webserver/
RUN go build -o webserver .

CMD ["./webserver"]
