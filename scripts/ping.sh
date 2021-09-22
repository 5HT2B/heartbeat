#!/usr/bin/env bash

# shellcheck disable=SC1091
source "$HOME/.env"
LOGGING_FILE="$HOME/.local/share/heartbeat.log"

if [[ -z "$(which xprintidle)" ]]; then
  echo "xprintidle not found, please install it!"
  exit 1
fi

# Check if kscreenlocker is running. Only works on KDE
SCREEN_LOCKED="$(pgrep kscreenlocker)"
# Check when the last keyboard or mouse event was sent
LAST_INPUT_MS="$(xprintidle)"

# Make sure the device was used in the last 2 minutes
# and make sure screen is unlocked
if [[ "$LAST_INPUT_MS" -lt 120000 && -z "$SCREEN_LOCKED" ]]; then
  {
    echo "$(date +"%Y/%m/%d %T") - Running Heartbeat"
    curl -s -X POST -H "Auth: $HEARTBEAT_AUTH" "$HEARTBEAT_HOSTNAME"
    echo ""
  } >> "$LOGGING_FILE" 2>&1
fi
