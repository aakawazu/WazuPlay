# FROM golang:1.14 AS binary

# WORKDIR /go/src/WazuPlay
# COPY . .

# RUN go build -o main /go/src/WazuPlay/cmd/oauth-server/main.go

# FROM alpine:latest
# WORKDIR /app
# EXPOSE 8080
# COPY --from=binary /go/src/WazuPlay/.env /app/
# COPY --from=binary /go/src/WazuPlay/main /app/cmd/oauth-server/

# CMD ["/app/cmd/oauth-server/main"]

FROM golang:1.14

WORKDIR /go/src/WazuPlay
COPY . .
WORKDIR /go/src/WazuPlay/cmd/oauth-server

RUN go build main.go

EXPOSE 8080

CMD ["/go/src/WazuPlay/cmd/oauth-server/main"]