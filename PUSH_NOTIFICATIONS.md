# Push Notifications for Heartbeat

## Overview

Server-initiated heartbeat checks allow the heartbeat server to actively monitor device status by sending push notifications that wake up PWA service workers and trigger immediate heartbeat responses. This works even when the PWA is completely closed.

## Setup Guide

### 1. Generate VAPID Keys

VAPID (Voluntary Application Server Identification) keys are required for web push notifications:

```bash
# Install web-push CLI globally
npm install -g web-push

# Generate VAPID key pair
web-push generate-vapid-keys

# Output example:
# Public Key: BEl62iUYgUivxIkv69yViEuiBIa-Ib9-SkvMeAtA3LFgDzkrxZJjSgSnfckjBJuBkr3qBUYIHBQFLXYp5Nksh8U
# Private Key: tUkzMcPbtl2-xZ4Z5A1OqiELPo2Pc9-SFNw6hU8_Hh0
```

- The private key should only be on your heartbeat server
- The public key needs to be in both server config AND PWA client code

### 2. Configure Heartbeat Server

Add the following to your `config/.env` file:

```bash
# VAPID Configuration for Push Notifications
VAPID_PUBLIC_KEY=BEl62iUYgUivxIkv69yViEuiBIa-Ib9-SkvMeAtA3LFgDzkrxZJjSgSnfckjBJuBkr3qBUYIHBQFLXYp5Nksh8U
VAPID_PRIVATE_KEY=tUkzMcPbtl2-xZ4Z5A1OqiELPo2Pc9-SFNw6hU8_Hh0
VAPID_SUBJECT=mailto:your-email@example.com

# Push Monitoring (Optional)
PUSH_CHECK_INTERVAL=5m          # How often to check for inactive devices  
PUSH_INACTIVITY_THRESHOLD=1h    # Trigger push if no heartbeat for this long
```

### 3. Configure PWA Client

The PWA now supports configuring the VAPID public key directly through the user interface:

1. **Open the PWA** in your browser or installed app
2. **Enter the VAPID Public Key** in the configuration section
3. **Save Configuration** to persist your settings
4. **Enable Server Checks** to subscribe to push notifications

The VAPID public key field accepts your generated public key (e.g., `BEl62iUYgUivxIkv69yViEuiBIa-Ib9-SkvMeAtA3LFgDzkrxZJjSgSnfckjBJuBkr3qBUYIHBQFLXYp5Nksh8U`)
