FROM golang:1.14-alpine3.11 AS binary

WORKDIR /go/src/WazuPlay
COPY . .
WORKDIR /go/src/WazuPlay/cmd/oauth-server

RUN go build main.go

FROM alpine:3.11
WORKDIR /go/src/WazuPlay/cmd/oauth-server
EXPOSE 8080

COPY --from=binary /go/src/WazuPlay/.env /go/src/WazuPlay
COPY --from=binary /go/src/WazuPlay/cmd/oauth-server/main /go/src/WazuPlay/cmd/oauth-server

CMD ["/go/src/WazuPlay/cmd/oauth-server/main"]