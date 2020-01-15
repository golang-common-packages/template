# build stage
FROM golang:alpine AS build-env
WORKDIR /

ENV GO111MODULE=on

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh gcc musl-dev && \
    apk --no-cache add ca-certificates

COPY ./ $GOPATH/src/github.com/golang-common-packages/template/
COPY ./config/main.yaml /config/main.yaml
COPY ./key/app.rsa /key/app.rsa
COPY ./key/app.rsa.pub /key/app.rsa.pub
COPY ./key/refresh.rsa /key/refresh.rsa
COPY ./key/refresh.rsa.pub /key/refresh.rsa.pub
RUN cd $GOPATH/src/github.com/golang-common-packages/template && \
    GOOS=linux GOARCH=amd64 go build -o /main

EXPOSE 3000
CMD ["/main"]