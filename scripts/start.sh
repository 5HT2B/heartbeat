#!/usr/bin/env bash

process="$(pgrep go)"

if [[ -z "$process" ]]; then
  cd /home/heartbeat/heartbeat
  exec /home/heartbeat/heartbeat/heartbeat -addr=YOUR_ADDR:YOUR_PORT -compress=true >> /home/heartbeat/heartbeat/heartbeat.log 2>&1
fi