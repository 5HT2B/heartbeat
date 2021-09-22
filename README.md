# [heartbeat](https://hb.l1v.in) 
[![time tracker](https://wakatime.com/badge/github/technically-functional/heartbeat.svg)](https://wakatime.com/badge/github/technically-functional/heartbeat)
[![Docker Pulls](https://img.shields.io/docker/pulls/technically-functional/heartbeat?logo=docker&logoColor=white)](https://hub.docker.com/r/technically-functional/heartbeat)
[![Docker Build](https://img.shields.io/github/workflow/status/technically-functional/heartbeat/docker-build?logo=docker&logoColor=white)](https://github.com/technically-functional/heartbeat/actions/workflows/docker-build.yml)
[![CodeFactor](https://img.shields.io/codefactor/grade/github/technically-functional/heartbeat?logo=codefactor&logoColor=white)](https://www.codefactor.io/repository/github/technically-functional/heartbeat)

A webpage to see when I was last active. Works by pinging the server from my computer or laptop every minute, as long as they have been used in the last minute.

For my laptop, this means if I have typed anything in the 2 minutes, for my phone it means if the screen was unlocked and on in the last 2 minutes.

This is my first time using Go, and I'm terrible at CSS, so this might not be using best practices. 

## Contributing

Contributions to fix my code are welcome, as well as any improvements.

To build:
```bash
git clone git@github.com:l1ving/heartbeat.git
cd heartbeat
make
```

To run:
```bash
# I recommend using genpasswd https://gist.github.com/l1ving/30f98284e9f92e1b47b4df6e05a063fc
AUTH='some secure token'
# We do not want to use echo because it appends a newline.
mkdir config
printf "$AUTH" > config/token

# Change the port to whatever you'd like. 
# Change localhost to your public IP if you'd like.
# Compression is optional, but enabled if not explicitly set.
./heartbeat -addr=localhost:8008 -compress=true
```

To test:

```bash
# Optionally add the -i flag before -X if you'd like more information.
curl -X POST -H "Auth: $AUTH" localhost:8008
```

or open localhost:8008 in a browser.

If you are having issues with a 403 even though you set it to POST and set your Auth header, PLEASE please make sure your `config/token` file does not have a trailing newline.

### Running server in production

I recommend using Caddy for automatic renewal + as a reverse proxy.
```
# Caddyfile example
hb.l1v.in {
  header Server Caddy "Heartbeat"
  reverse_proxy localhost:6060
}
```

There is also a docker image available with the following command, or checkout the
[`update.sh`](https://github.com/l1ving/heartbeat/blob/master/scripts/update.sh) script for automatically
updating a live docker image.
```bash
docker pull l1ving/heartbeat:latest
```

### Running client on Linux

Copy systemd service (`-client`) files to `~/.config/systemd/user/` and edit the `ExecStart` accordingly.
Make sure the path matches the *full* path to your `ping.sh`.

Also add this to the end of your `$HOME/.env` to include the required env variables.

```bash
# Note: Please use your own token. This is simpily an example.
export HEARTBEAT_AUTH='Cn$Sn61rt6knaSU06NEntzVTMrLnBN&c15UBbdkn6;vJzJ9D' # Single quotes to avoid escaping issues.
export HEARTBEAT_HOSTNAME="localhost:8008" # Change your IP to your www IP.
```

Then
```bash
mkdir -p "$HOME/.local/share" # Create default logging folder. Edit `ping.sh` if you don't like this.
# Install xprintidle with your distro's package manager
systemctl --user enable heartbeat-client.timer
systemctl --user start heartbeat-client.timer
```

### Running client on Android (tasker)

You should be able to import each of these tasker profiles from the [tasker folder](https://github.com/l1ving/heartbeat/tree/master/tasker). Make sure each of them are enabled, and **edit the server address and Auth token inside `Ping.tsk.xml` before importing**.

Make sure to allow running in background, and all of the other optimizations tasker recommends to be sure that it runs. If there is any issue importing it, please [make an issue](https://github.com/l1ving/heartbeat/issues/new) in this repo.

You will also need to create your own profile that runs every 2 minutes, to run the "ping task", and make sure the Display Off and Display Unlocked profiles are enabled.

### Running client on anything else

I wanted something as simple to interface with as possible, which is why it's relatively easy to get it working on Android as well.

Anything that can run curl can update the last beat. 

I **highly** recommend checking the last time an input device such as a keyboard or mouse was used, much like `ping.sh` does.

```bash
curl -s -X POST -H "Auth: $HEARTBEAT_AUTH" $HEARTBEAT_HOSTNAME
```
