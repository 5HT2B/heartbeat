#!/bin/bash

docker pull l1ving/heartbeat:latest

if [[ "$1" != "FIRST_RUN" ]]; then
  CONTAINER_ID="$(docker ps -f name=heartbeat --format "{{.ID}}" | head -n 1)"
  echo "Stopping container $CONTAINER_ID"
  docker stop "$CONTAINER_ID"
  docker rm "$CONTAINER_ID"
else
  echo "MAKE SURE TO CREATE 'token' inside \''$HOME'/heartbeat\'"
fi

docker run --name heartbeat \
  -e ADDRESS="localhost:6011" \
  --mount type=bind,source="$HOME"/heartbeat,target=/heartbeat/config \
  --network host -d \
  l1ving/heartbeat
