#!/usr/bin/env bash

source ~/.profile
LOGGING_FILE="$HOME/.local/share/heartbeat.log"

if [[ -z "$(which xprintidle)" ]]; then
  echo "xprintidle not found, please install it!"
  exit 1
fi

# Allow for 2 minutes per cronjob buffer
if [[ "$(xprintidle)" -lt 120000 ]]; then
  echo "$(date +"%Y/%m/%d %T") - Running Heartbeat" >> "$LOGGING_FILE"
  curl -s -X POST -H "Auth: $HEARTBEAT_AUTH" "$HEARTBEAT_HOSTNAME" >> "$LOGGING_FILE" 2>&1
  echo "" >> "$LOGGING_FILE"
fi
