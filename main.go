//go:generate qtc -dir=templates
package main

import (
	"flag"
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"
	"log"
	"os"
	"time"
)

//goland:noinspection HttpUrlsUsage
var (
	debug         = flag.Bool("debug", false, "Whether to print debug output")
	authTokenFlag = flag.String("token", "", "An alternative token to be used when debugging")
	protocol      = "http://"             // set in .env
	addr          = "0.0.0.0:6060"        // set in .env
	serverName    = "A Development Build" // set in .env
	authToken     = ""                    // set in .env // TODO: Add support for multi tokens
	gitCommitHash = "<unknown>"           // This is changed with compile flags in Makefile
	timeFormat    = "Jan 02 2006 15:04:05 MST"
	uptimeTimer   = time.Now().Unix()
)

func main() {
	flag.Parse()
	setupEnv()
	rdb, rjh = SetupDatabase()
	go SetupLocalValues()
	go SetupDatabaseSaving()
	log.Printf("- Running heartbeat on " + protocol + addr)

	h := RequestHandler
	h = fasthttp.CompressHandler(h)

	if err := fasthttp.ListenAndServe(addr, h); err != nil {
		log.Fatalf("- Error in ListenAndServe: %s", err)
	}

	defer func() {
		if err := rdb.Close(); err != nil {
			log.Fatalf("goredis - failed to communicate to redis-server: %v", err)
		}
	}()
}

func setupEnv() {
	err := godotenv.Load("config/.env")
	if err != nil {
		log.Fatalf("Fatal error: %v", err)
		return
	}

	if *debug && len(*authTokenFlag) > 0 {
		authToken = *authTokenFlag
	} else {
		if ev := os.Getenv("HB_TOKEN"); len(ev) == 0 {
			log.Fatalf("HB_TOKEN not set in config/.env")
			return
		} else {
			authToken = ev
		}
	}

	if ev := os.Getenv("HB_PROTOCOL"); len(ev) > 0 {
		protocol = ev
	}
	if ev := os.Getenv("HB_ADDR"); len(ev) > 0 {
		addr = ev
	}
	if ev := os.Getenv("HB_SERVER_NAME"); len(ev) > 0 {
		serverName = ev
	}
	if ev := os.Getenv("HB_GITHUB_LINK"); len(ev) > 0 {
		gitRepo = ev
	}
	if ev := os.Getenv("REDIS_ADDR"); len(ev) > 0 {
		redisAddr = ev
	}
	if ev := os.Getenv("REDIS_PASS"); len(ev) > 0 {
		redisPass = ev
	}
}
