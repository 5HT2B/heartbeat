# [heartbeat](https://hb.l1v.in) 
[![time tracker](https://wakatime.com/badge/github/l1ving/heartbeat.svg)](https://wakatime.com/badge/github/l1ving/heartbeat)
[![Docker Pulls](https://img.shields.io/docker/pulls/l1ving/heartbeat?logo=docker&logoColor=white)](https://hub.docker.com/r/l1ving/heartbeat)
[![Docker Build](https://img.shields.io/github/workflow/status/technically-functional/heartbeat/docker-build?logo=docker&logoColor=white)](https://github.com/technically-functional/heartbeat/actions/workflows/docker-build.yml)
[![CodeFactor](https://img.shields.io/codefactor/grade/github/technically-functional/heartbeat?logo=codefactor&logoColor=white)](https://www.codefactor.io/repository/github/technically-functional/heartbeat)

A service to see when a device was last active. Works by pinging the server every minute, from any device, as long as said device is unlocked and being used (ie, you typed or used the mouse in the last two minutes).

## Contributing

Contributions to fix code are welcome, as are any improvements.

## Usage

```bash
git clone git@github.com:technically-functional/heartbeat.git
cd heartbeat
echo "HB_TOKEN=authenticationTokenMakeThisSecure" >> config/.env
HB_PATH=$(pwd)/config docker-compose up --build
```

<!-- OLD METHOD 
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
-->
To test a ping locally:

```bash
# Optionally add the -i flag if you'd like more information.
curl -X POST -H "Auth: $AUTH" -H "Device: laptop" localhost:6060/api/beat

# Optionally, you can set the token with an auth flag instead of .env for debugging
./heartbeat -debug -token some_token_here
```

or open localhost:6060 in a browser to view the webpage.

## Running server or client in production

See [`DOCS.md#running-server-in-production`](https://github.com/technically-functional/heartbeat/blob/master/DOCS.md#running-server-in-production) for more information.

## API

See [`DOCS.md#api`](https://github.com/technically-functional/heartbeat/blob/master/DOCS.md#api) for more information.
