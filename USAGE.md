# Usage

## With Docker

This requires `docker-compose` (1.29.2 or newer), `docker` and `git`.

```bash
git clone git@github.com:technically-functional/heartbeat.git
cd heartbeat
echo "HB_TOKEN=authenticationTokenMakeThisSecure" >> config/.env
export HB_PATH=$(pwd)/config
# Add --build to build from your local files instead of using the pre-built image 
docker-compose up
```

## Without Docker

This requires `go` (1.16 or newer), `redis-server` (6.2.5 or newer), [RedisJSON](https://github.com/RedisJSON/RedisJSON) and `git`.

```bash
# Run redis-server with the following command in the background, or in another window
redis-server path/to/hb/config/redis.conf --loadmodule path/to/RedisJSON/target/release/librejson.so

# Run the following in a new terminal
git clone git@github.com:technically-functional/heartbeat.git
cd heartbeat
# Make sure to edit REDIS_ADDR inside config/.env to localhost:6379
echo "HB_TOKEN=authenticationTokenMakeThisSecure" >> config/.env
make # Build the binary
./heartbeat
```

## Testing

To test a ping locally:

```bash
# Optionally add the -i flag if you'd like more information.
curl -X POST -H "Auth: $AUTH" -H "Device: laptop" localhost:6060/api/beat

# Optionally, you can set the token with an auth flag instead of .env for debugging
./heartbeat -debug -token some_token_here
```

or open localhost:6060 in a browser to view the webpage.

## Debugging

- Can't connect using Docker?

    The default port is `6060`, and you should be able to access `localhost:6060`. This is set in `config/.env` and `docker-compose.yml`.

    If you are unable to connect from localhost, make sure these are set to your desired port, and check the `docker-compose` log for issues.

- Can't connect without Docker?

    The default port is `6060`, set in `config/.env` with `HB_ADDR`. If `./heartbeat` isn't throwing any errors, please check that you have the right port.

- Can't `POST` to `/api/beat`?

    Try running `./heartbeat -debug -token some_token_here`, which will override the default token, to help debug the issue.

- Heartbeat can't read the `config/.env` when using Docker?

    Make sure `export HB_PATH` is pointing to the config folder inside your Heartbeat folder.

- `dial tcp: lookup database: no such host`

    Heartbeat can't connect to the Redis database. If you're not using the Docker image, make sure that you ran `redis-server` before `./heartbeat`.
