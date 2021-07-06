#!/bin/bash

docker pull l1ving/heartbeat:latest
CONTAINER_ID="$(docker ps -f name=heartbeat --format "{{.ID}}" | head -n 1)"

echo "Stopping container $CONTAINER_ID"
docker stop "$CONTAINER_ID"
docker rm "$CONTAINER_ID"

docker run --name heartbeat \
  -e ADDRESS="localhost:6011" \
  --mount type=bind,source=/home/liv/heartbeat,target=/heartbeat-files \
  --network host -d \
  l1ving/heartbeat
