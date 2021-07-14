FROM golang:1.16.6-buster
RUN apt update
RUN apt install -y git
WORKDIR /root/
RUN git clone https://github.com/panzerox123/blockcert
WORKDIR /root/blockcert
RUN go build blockcert.go
RUN go install
VOLUME [ "/data" ]
WORKDIR /data
EXPOSE 8080
CMD ["/go/bin/blockcert"]

