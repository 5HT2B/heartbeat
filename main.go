//go:generate qtc -dir=templates
package main

import (
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"
	"log"
	"os"
)

//goland:noinspection HttpUrlsUsage
var (
	protocol                          = "http://"             // set in .env
	addr                              = "localhost:6060"      // set in .env
	serverName                        = "A Development Build" // set in .env
	authToken                         = []byte(ReadFileUnsafe("config/token"))
	gitCommitHash                     = "<unknown>" // This is changed with compile flags in Makefile
	timeFormat                        = "Jan 02 2006 15:04:05 MST"
	lastBeat, missingBeat, totalBeats = ReadLastBeatSafe()
	totalVisits                       = ReadGetRequestsSafe()
)

func main() {
	setupEnv()
	log.Printf("- Running heartbeat on " + protocol + addr)

	h := RequestHandler
	h = fasthttp.CompressHandler(h)

	if err := fasthttp.ListenAndServe(addr, h); err != nil {
		log.Fatalf("- Error in ListenAndServe: %s", err)
	}
}

func setupEnv() {
	err := godotenv.Load("config/.env")
	if err != nil {
		log.Fatalf("Fatal error: %v", err)
		return
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
}
