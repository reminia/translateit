FROM golang:1.20-alpine3.17

WORKDIR /opt
RUN mkdir translateit
WORKDIR /opt/translateit
COPY *.go  .
COPY *.mod  .
COPY *.sum  .
COPY Makefile .
RUN apk add make
RUN make build

EXPOSE 8080

CMD ["./translate"]