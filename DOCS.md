# Docs

A few technical aspects are kept documented here for future reference.
For the most part, functions will be commented.

## Running server in production

Using Caddy or Nginx + Certbot for automatic renewal and reverse proxying is recommended.
```
# Caddyfile example
hb.l1v.in {
  header Server Caddy "Heartbeat"
  reverse_proxy localhost:6060
}
```

```
# Nginx example
server {
    server_name hb.l1v.in;
    location / {
        proxy_pass localhost:6060;
        proxy_set_header X-Real-IP $remote_addr;
    }
    # Automatically managed certbot stuff goes here
} 
```

There is also a docker image available with the following command, or checkout the
[`update.sh`](https://github.com/technically-functional/heartbeat/blob/master/scripts/update.sh) script for automatically updating a live docker image.
```bash
# Simply pull the image
docker pull l1ving/heartbeat:latest
# Run the service under docker. Do not use FIRST_RUN if you have run it before.
./update.sh FIRST_RUN
```

### Running client on Android (tasker)

You should be able to import each of these tasker profiles from the [tasker folder](https://github.com/technically-functional/heartbeat/tree/master/tasker). Make sure each of them are enabled, and **edit the server address and Auth token inside `Ping.tsk.xml` before importing**.

Make sure to allow running in background, as well as the other optimizations which Tasker recommends, in order to be sure that it runs.
If there is any issue importing it, please [make an issue](https://github.com/technically-functional/heartbeat/issues/new) in this repo.

You will also need to create your own profile that runs every 2 minutes, to run the "ping task", and make sure the Display Off and Display Unlocked profiles are enabled.

### Running client on anything else

Please see [`heartbeat-unix`](https://github.com/technically-functional/heartbeat-unix) for a (mostly universal) example that will run on almost any \*NIX system.

## API

Available endpoints:
- `/api/beat` (POST)
  - `curl -X POST -H "Auth: $AUTH" -H "Device: $DEVICE_NAME localhost:6060/api/beat`
  - Requires the Auth and Device headers to be set
- `/api/update/stats` (POST)
  - Requires the Auth header to be set
  - `curl -X POST -H "Auth: $AUTH" localhost:6060/api/update/stats -d '$STATS_JSON'`
- `/api/update/devices` (POST)
  - Requires the Auth header to be set
  - `curl -X POST -H "Auth: $AUTH" localhost:6060/api/update/devices -d '$DEVICES_JSON'`
- `/api/stats` (GET)
- `/api/devices` (GET)

## Redis JSONPath structure example

```bash
.
├── beats
│   ├── <HeartbeatBeat>        # {"device_name": "laptop", "timestamp": 1632748096}
│   ├── <HeartbeatBeat>        # {"device_name": "phone", "timestamp": 1632748137}
│   └── <HeartbeatBeat>        # {"device_name": "laptop", "timestamp": 1632748682}
│
├── devices
│   ├── <HeartbeatDevice>      # {"device_name": "laptop", "last_beat": <HeartbeatBeat>, "total_beats": 12903, "longest_missing_beat", 1201}
│   └── <HeartbeatDevice>      # {"device_name": "phone", "last_beat": <HeartbeatBeat>, "total_beats": 1952, "longest_missing_beat", 3219}
│
├── last_beat <HeartbeatBeat>  # {"device_name": "laptop", "timestamp": 1632748096}
└── stats <HeartbeatStats>     # {"total_visits": 45012, "total_uptime": 892340, "total_beats": 14855, "longest_missing_beat", 3219}
```
