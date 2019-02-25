FROM golang:alpine AS builder

ADD . /go/src/go_base
WORKDIR /go/src/go_base

# install git
RUN apk add --no-cache git bzr
RUN go get golang.org/x/sys/unix
RUN go get -u -v github.com/kardianos/govendor && \
   govendor sync
RUN GOOS=linux GOARCH=amd64 go build -v -o /go/src/go_base/main

FROM alpine
WORKDIR /root
RUN apk add -U tzdata && \
   ln -sf /usr/share/zoneinfo/Asia/Shanghai  /etc/localtime
COPY --from=builder /go/src/go_base .
#EXPOSE 8080
CMD [ "./main" ]
