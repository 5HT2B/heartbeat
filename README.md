# heartbeat

A webpage to see when I was last active. Works by pinging the server from my computer or laptop every minute, as long as they have been used in the last minute.

For my laptop, this means if I have typed anything in the last minute, for my phone it means if the screen was on.

This is my first time using Go and I'm terrible at CSS, so this might not be using best practices. 

## Contributing

Contributions to fix my code are welcome.
I also didn't bother googling how to implement a favicon, so we're just living without it.

To build:
```bash
git clone git@github.com:l1ving/heartbeat.git
cd heartbeat
make
```

To run:
```bash
# I recommend using genpasswd https://gist.github.com/l1ving/30f98284e9f92e1b47b4df6e05a063fc
AUTH="some secure token"
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
