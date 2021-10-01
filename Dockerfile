FROM golang:1.16.5
ARG COMMIT="latest"

RUN mkdir -p /heartbeat/config
ADD . /heartbeat
WORKDIR /heartbeat

ENV PATH "${PATH}:${GOPATH}/bin"
RUN make update build

CMD /heartbeat/heartbeat
