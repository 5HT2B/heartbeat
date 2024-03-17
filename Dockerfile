FROM golang:1.20.14
ARG COMMIT="latest"

RUN mkdir -p /heartbeat/config
ADD . /heartbeat
WORKDIR /heartbeat

ENV PATH "${PATH}:${GOPATH}/bin"
RUN make deps build

CMD /heartbeat/heartbeat
