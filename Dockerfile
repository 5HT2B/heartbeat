FROM golang:1.16.5
ARG COMMIT="latest"

RUN mkdir /heartbeat
RUN mkdir /heartbeat-files
ADD . /heartbeat
WORKDIR /heartbeat

RUN go build -ldflags "-X main.gitCommitHash=${COMMIT}" -o heartbeat .

ENV ADDRESS "localhost:6060"

ADD www /heartbeat-files/www
WORKDIR /heartbeat-files
CMD /heartbeat/heartbeat -addr ${ADDRESS}
