# [heartbeat](https://hb.l1v.in)

A webpage to see when I was last active. Works by pinging the server from my computer or laptop every minute, as long as they have been used in the last minute.

For my laptop, this means if I have typed anything in the 2 minutes, for my phone it means if the screen was unlocked and on in the last 2 minutes.

This is my first time using Go and I'm terrible at CSS, so this might not be using best practices. 

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
printf "$AUTH" > token

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

If you are having issues with a 403 even though you set it to POST and set your Auth header, PLEASE please make sure your `token` file does not have a trailing newline.

### Running server in production

Follow the above instructions, just skip the testing section.

Copy systemd service (`-server`) files to `~/.config/systemd/user/` and edit the `ExecStart` accordingly. 
Make sure the path matches the *full* path to your `start.sh`.

Also edit your `start.sh` to match the path to your `heartbeat` binary.

Then 
```bash
systemctl --user enable heartbeat-server.timer
systemctl --user start heartbeat-server.timer
```

This will run the service automatically, and restart it if it dies.

You can update Heartbeat on the server by simpily 

```bash
cd ~/path/to/heartbeat
git pull
make
systemctl --user stop heartbeat-server.service
# Automatically restarts with the new binary
```

If you'd like to run Heartbeat on **ports 80 or 443**, please read [this](https://superuser.com/a/892391).

### Running client on Linux

Copy systemd service (`-client`) files to `~/.config/systemd/user/` and edit the `ExecStart` accordingly.
Make sure the path matches the *full* path to your `ping.sh`.

Also add this to the end of your `~/.profile` to include the required env variables.

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

You *should* be able to import it from [this](https://taskernet.com/shares/?user=AS35m8nEM1zYe7Hwhnr%2FmqY6FLKRigBn5KsRZjpRBVQ5kVsRat6L8dgyXksiNLNHQ5ycPrAdpiCS860%3D&id=Profile%3ADisplay+Unlocked) link on your Android.

Make sure to allow running in background, and all of the other optimizations tasker recommends to be sure that it runs. If there is any issue importing it, please [make an issue](https://github.com/l1ving/heartbeat/issues/new) in this repo.

Also make sure to edit the auth header in the `ping` task.

### Running client on anything else

I wanted something as simple to interface with as possible, which is why it's relatively easy to get it working on Android as well.

Anything that can run curl can update the last beat. 

I **highly** recommend checking the last time an input device such as a keyboard or mouse was used, much like `ping.sh` does.

```bash
curl -s -X POST -H "Auth: $HEARTBEAT_AUTH" $HEARTBEAT_HOSTNAME
```
