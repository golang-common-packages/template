# Builder
FROM golang:1.15.2-alpine3.12 as builder

RUN apk update && apk upgrade && \
    apk --update add git make

WORKDIR /app

COPY . .

RUN make engine

# Distribution
FROM alpine:latest

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    mkdir /app 

WORKDIR /app 

EXPOSE 4000

COPY --from=builder /app/engine /app

COPY --from=builder /app/config.json /app

CMD /app/engine