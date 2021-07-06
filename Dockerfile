FROM golang:1.16.5

RUN mkdir /heartbeat
RUN mkdir /heartbeat-files
ADD . /heartbeat
WORKDIR /heartbeat

RUN go build -o heartbeat .

ENV ADDRESS "localhost:6060"

WORKDIR /heartbeat-files
CMD /heartbeat/heartbeat -addr $ADDRESS
