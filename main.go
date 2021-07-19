package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/Ferluci/fast-realip"
	"github.com/valyala/fasthttp"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	addr                              = flag.String("addr", "localhost:6060", "TCP address to listen to")
	compress                          = flag.Bool("compress", true, "Whether to enable transparent response compression")
	useTls                            = flag.Bool("tls", false, "Whether to enable TLS")
	tlsCert                           = flag.String("cert", "", "Full certificate file path")
	tlsKey                            = flag.String("key", "", "Full key file path")
	unknown403                        = flag.Bool("unknown403", true, "Return 403 on unknown paths.")
	serverName                        = flag.String("name", "Liv's Heartbeat", "The name of the server to use")
	authToken                         = []byte(ReadFileUnsafe("config/token"))
	pathRoot                          = []byte("/")
	pathFavicon                       = []byte("/favicon.ico")
	gitCommitHash                     = "A Development Version" // This is changed with compile flags in Makefile
	timeFormat                        = "Jan 02 2006 15:04:05 MST"
	htmlFile                          = ReadFileUnsafe("www/index.html")
	lastBeat, missingBeat, totalBeats = ReadLastBeatSafe()
	totalVisits                       = ReadGetRequestsSafe()
)

func main() {
	flag.Parse()

	protocol := "http"
	if *useTls {
		protocol += "s"
	}

	log.Printf("- Running heartbeat on " + protocol + "://" + *addr)

	h := RequestHandler
	if *compress {
		h = fasthttp.CompressHandler(h)
	}

	if *useTls && len(*tlsCert) > 0 && len(*tlsKey) > 0 {
		if err := fasthttp.ListenAndServeTLS(*addr, *tlsCert, *tlsKey, h); err != nil {
			log.Fatalf("- Error in ListenAndServeTLS: %s", err)
		}
	} else {
		if err := fasthttp.ListenAndServe(*addr, h); err != nil {
			log.Fatalf("- Error in ListenAndServe: %s", err)
		}
	}
}

func RequestHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set(fasthttp.HeaderServer, *serverName)

	path := ctx.Path()
	requestPath := string(path)

	if bytes.Equal(path, pathRoot) {
		HandleRoot(ctx)
		return
	}

	if !ctx.IsGet() { // Don't allow get on non-root
		HandleUnknown(ctx)
		return
	}

	if bytes.Equal(path, pathFavicon) {
		HandleFavicon(ctx)
		return
	}

	pathFile := ""
	if requestPath == "/stats" {
		pathFile = "config/get_requests"
	} else {
		pathFile = "www" + requestPath // sep is not / because requestPath is prefixed with a /
	}

	err := ServeFile(ctx, pathFile, false)

	// Try to find corresponding file with the .html suffix added
	if err != nil {
		pathFile = pathFile + ".html"

		err = ServeFile(ctx, pathFile, true)

		if err != nil {
			return
		}
	}

	totalVisits++
}

func ServeFile(ctx *fasthttp.RequestCtx, file string, handleErr bool) error {
	content, err := ReadFile(file)

	if err == nil {
		switch {
		case strings.HasSuffix(file, ".css"):
			ctx.Response.Header.Set(fasthttp.HeaderContentType, "text/css; charset=utf-8")
		case strings.HasSuffix(file, ".html"):
			ctx.Response.Header.Set(fasthttp.HeaderContentType, "text/html; charset=utf-8")
		}

		fmt.Fprint(ctx, content)
	} else if handleErr {
		HandleUnknown(ctx)
	}

	return err
}

func HandleRoot(ctx *fasthttp.RequestCtx) {
	// Serve the HTML page if it is a GET request
	if ctx.IsGet() {
		ctx.Response.Header.Set(fasthttp.HeaderContentType, "text/html; charset=utf-8")
		fmt.Fprint(ctx, GetHtml())
		totalVisits++
		return
	}

	// We only want to allow POST requests if we're not serving a page
	if !ctx.IsPost() {
		ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
		fmt.Fprint(ctx, "400 Bad Request\n")
		log.Printf("- Returned 400 to %s - tried to connect with '%s'", realip.FromRequest(ctx), ctx.Method())
		return
	}

	// The authentication key provided with said Auth header
	header := ctx.Request.Header.Peek("Auth")

	// Make sure Auth key is correct
	if !bytes.Equal(header, authToken) {
		HandleUnknown(ctx)
		return
	}

	// Now that we know it is a POST and the Auth key is correct, return the unix timestamp and write it to a file
	HandleSuccessfulBeat(ctx)
}

func HandleUnknown(ctx *fasthttp.RequestCtx) {
	if *unknown403 {
		ctx.Response.SetStatusCode(fasthttp.StatusForbidden)
		fmt.Fprint(ctx, "403 Forbidden\n")
		log.Printf("- Returned 403 to %s - tried to connect with '%s' to '%s'", realip.FromRequest(ctx), ctx.Request.Header.Peek("Auth"), ctx.Path())
	}
}

func HandleFavicon(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set(fasthttp.HeaderContentType, "image/x-icon")
	f, _ := os.Open("www/favicon.ico")
	defer f.Close()
	io.Copy(ctx, f)
}

func GetHtml() string {
	currentTime := time.Now()
	currentBeatDifference := currentTime.Unix() - lastBeat

	// We also want to live update the current difference, instead of only when receiving a beat.
	if currentBeatDifference > missingBeat {
		missingBeat = currentBeatDifference
	}

	lastSeen := time.Unix(lastBeat, 0).Format(timeFormat)
	timeDifference := TimeDifference(lastBeat, currentTime)
	missingBeatFmt := FormattedTime(missingBeat)
	totalBeatsFmt := FormattedNum(totalBeats)
	currentTimeStr := currentTime.Format(timeFormat)

	htmlTemp := htmlFile
	htmlTemp = strings.Replace(htmlTemp, "LAST_SEEN", lastSeen, 2)
	htmlTemp = strings.Replace(htmlTemp, "RELATIVE_TIME", timeDifference, 1)
	htmlTemp = strings.Replace(htmlTemp, "LONGEST_ABSENCE", missingBeatFmt, 1)
	htmlTemp = strings.Replace(htmlTemp, "TOTAL_BEATS", totalBeatsFmt, 1)
	htmlTemp = strings.Replace(htmlTemp, "CURRENT_TIME", currentTimeStr, 1)
	htmlTemp = strings.Replace(htmlTemp, "GIT_HASH", gitCommitHash, 2)
	htmlTemp = strings.Replace(htmlTemp, "GIT_REPO", "https://github.com/l1ving/heartbeat", 2)
	htmlTemp = strings.Replace(htmlTemp, "SERVER_NAME", *serverName, 3)

	return htmlTemp
}

func HandleSuccessfulBeat(ctx *fasthttp.RequestCtx) {
	totalBeats += 1
	newLastBeat := time.Now().Unix()
	currentBeatDifference := newLastBeat - lastBeat

	if currentBeatDifference > missingBeat {
		missingBeat = currentBeatDifference
	}

	lastBeatStr := strconv.FormatInt(newLastBeat, 10)
	missingBeatStr := strconv.FormatInt(missingBeat, 10)
	totalBeatsStr := strconv.FormatInt(totalBeats, 10)

	fmt.Fprintf(ctx, "%v\n", lastBeatStr)
	log.Printf("- Successful beat from %s", realip.FromRequest(ctx))

	lastBeat = newLastBeat
	WriteToFile("config/last_beat", lastBeatStr+":"+missingBeatStr+":"+totalBeatsStr)
	WriteGetRequestsFile(totalVisits)
}
