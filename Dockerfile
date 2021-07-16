FROM golang:1.16.5
ARG COMMIT="latest"

RUN mkdir /heartbeat \
 && mkdir /heartbeat/config
ADD . /heartbeat
WORKDIR /heartbeat

RUN go build -ldflags "-X main.gitCommitHash=${COMMIT}" -o heartbeat .

ENV ADDRESS "localhost:6060"
CMD /heartbeat/heartbeat -addr ${ADDRESS}
