#!/bin/bash

source "$HOME/.profile"
if [[ -z "$HB_PATH" ]]; then
  echo "HB_PATH not set!"
  exit 1
fi

docker pull l1ving/heartbeat:latest

if [[ "$1" != "FIRST_RUN" ]]; then
  docker stop heartbeat || echo "Could not stop missing container heartbeat"
  docker rm heartbeat || echo "Could not remove missing container heartbeat"
else
  echo "MAKE SURE TO CREATE 'token' inside '$FOH_PATH'"
fi

docker run --name heartbeat \
  -e ADDRESS="localhost:6011" \
  --mount type=bind,source="$FOH_PATH",target=/heartbeat/config \
  --network host -d \
  l1ving/heartbeat
