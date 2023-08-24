FROM golang:1.17.6 as build
ARG COMMIT="latest"
WORKDIR /heartbeat
RUN mkdir -p /heartbeat/config
COPY . .

ENV PATH "${PATH}:${GOPATH}/bin"
RUN make deps build

FROM gcr.io/distroless/base-debian11
COPY --from=build /heartbeat/heartbeat /heartbeat/heartbeat
COPY --from=build /heartbeat/www /heartbeat/www
WORKDIR /heartbeat
ENTRYPOINT [ "/heartbeat/heartbeat" ]
