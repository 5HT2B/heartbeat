# Usage

**Before following either of the following run examples**, rename `config/.env.example` to `config/.env` and
`config/redis.conf.example` to `config/redis.conf`.

To change the port Heartbeat is running on, regardless of running method, edit `config/.env` and change `HB_PORT` and `HB_ADDR`
(as well as `docker-compose.yml` if using Docker).

## With Docker

This requires `docker-compose` (1.29.2 or newer), `docker` and `git`.

```bash
git clone git@github.com:5HT2B/heartbeat.git
cd heartbeat
echo "HB_TOKEN=authenticationTokenMakeThisSecure" >> config/.env
# Add --build to build from your local files instead of using the pre-built image 
docker-compose up
```

## Without Docker

This requires `go` (1.20 or newer), `redis-server` (6.2.5 or newer), [RedisJSON](https://github.com/RedisJSON/RedisJSON) and `git`.

```bash
# Make sure to edit dir inside config/redis.conf to ./config
# Make sure to edit REDIS_ADDR inside config/.env to localhost:6379

# Run redis-server with the following command in the background, or in another window
redis-server config/redis.conf --loadmodule /path/to/RedisJSON/target/release/librejson.so

# Run the following in a new terminal
git clone git@github.com:5HT2B/heartbeat.git
cd heartbeat
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

## Webhooks

Set `HB_WEBHOOK_URL`, `HB_WEBHOOK_LEVEL` and `HB_LIVE_URL` in `config/.env` to enable webhooks, see `config/.env.example` for more info.

For `HB_WEBHOOK_LEVEL`, see `WebhookLevel` in `webhook.go` for an explanation.

## Debugging

- Can't connect using Docker?

  The default port is `6060`, and you should be able to access `localhost:6060`. This is set in `config/.env`.

  If you are unable to connect from localhost, make sure these are set to your desired port, and check the `docker-compose logs` for issues.

- Can't connect without Docker?

  The default port is `6060`, set in `config/.env` with `HB_ADDR` and `HB_PORT`. If `./heartbeat` isn't throwing any errors, please check
  that you have the right port.

- Can't `POST` to `/api/beat`?

  Try running `./heartbeat -debug -token some_token_here`, which will override the default token, to help debug the issue.

- Heartbeat can't read the `config/.env` when using Docker?

  Make sure you are editing the `config/.env` which is inside the same folder as your `docker-compose.yml`.

- `dial tcp: lookup database: no such host`

  Heartbeat can't connect to the Redis database. If you're not using the Docker image, make sure that you ran `redis-server`
  before `./heartbeat`.
