FROM golang:1.20-alpine3.17

WORKDIR /opt
RUN mkdir translateit
WORKDIR /opt/translateit
COPY *.go  .
COPY *.mod  .
COPY *.sum  .
RUN go build translate.go specs.go
EXPOSE 8080

CMD ["./translate"]