FROM golang:1.14-alpine3.11 AS binary

WORKDIR /go/src/WazuPlay
COPY . .
WORKDIR /go/src/WazuPlay/cmd/audio-server
RUN go build main.go

FROM alfg/ffmpeg:latest
WORKDIR /go/src/WazuPlay/cmd/audio-server
EXPOSE 8080

COPY --from=binary /go/src/WazuPlay/.env /go/src/WazuPlay
COPY --from=binary /go/src/WazuPlay/cmd/audio-server/main /go/src/WazuPlay/cmd/audio-server
RUN mkdir -p /wazuplay-files/audio/tmp

CMD ["/go/src/WazuPlay/cmd/audio-server/main"]