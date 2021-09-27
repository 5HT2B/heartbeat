# [heartbeat](https://hb.l1v.in) 
[![time tracker](https://wakatime.com/badge/github/l1ving/heartbeat.svg)](https://wakatime.com/badge/github/l1ving/heartbeat)
[![Docker Pulls](https://img.shields.io/docker/pulls/l1ving/heartbeat?logo=docker&logoColor=white)](https://hub.docker.com/r/l1ving/heartbeat)
[![Docker Build](https://img.shields.io/github/workflow/status/technically-functional/heartbeat/docker-build?logo=docker&logoColor=white)](https://github.com/technically-functional/heartbeat/actions/workflows/docker-build.yml)
[![CodeFactor](https://img.shields.io/codefactor/grade/github/technically-functional/heartbeat?logo=codefactor&logoColor=white)](https://www.codefactor.io/repository/github/technically-functional/heartbeat)

A service to see when a device was last active. Works by pinging the server every minute, from any device, as long as said device is unlocked and being used (ie, you typed or used the mouse in the last two minutes).

## Contributing

Contributions to fix code are welcome, as are any improvements.

To build:
```bash
git clone git@github.com:technically-functional/heartbeat.git
cd heartbeat
make
```

To run:
```bash
# Use genpasswd to create a token, or another random password generator
# https://gist.github.com/l1ving/30f98284e9f92e1b47b4df6e05a063fc

edit config/.env
# And set HB_TOKEN to a secure token
# Change the port to whatever you'd like
# Changing localhost to a public IP isn't recommended without setting up https
# Ideally, you could also use a reverse proxy on localhost + certbot

# Run heartbeat now that your config/.env is setup
./heartbeat
```

To test a ping locally:

```bash
# Optionally add the -i flag if you'd like more information.
curl -X POST -H "Auth: $AUTH" localhost:8008/api/beat
```

or open localhost:6060 in a browser to view the webpage.

### Running server in production

I recommend using Caddy or Nginx + Certbot for automatic renewal and reverse proxying.
```
# Caddyfile example
hb.l1v.in {
  header Server Caddy "Heartbeat"
  reverse_proxy localhost:6060
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

You should be able to import each of these tasker profiles from the [tasker folder](https://github.com/l1ving/heartbeat/tree/master/tasker). Make sure each of them are enabled, and **edit the server address and Auth token inside `Ping.tsk.xml` before importing**.

Make sure to allow running in background, and all of the other optimizations tasker recommends to be sure that it runs. If there is any issue importing it, please [make an issue](https://github.com/l1ving/heartbeat/issues/new) in this repo.

You will also need to create your own profile that runs every 2 minutes, to run the "ping task", and make sure the Display Off and Display Unlocked profiles are enabled.

### Running client on anything else

Please see [`heartbeat-unix`](https://github.com/technically-functional/heartbeat-unix) for a (mostly universal) example that will run on almost any \*NIX system.
