FROM golang:1.14

WORKDIR /go/src/app
COPY ./oauth-server/. .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["go", "run", "/go/src/app/."]