package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/valyala/fasthttp"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	addr                  = flag.String("addr", ":6060", "TCP address to listen to")
	compress              = flag.Bool("compress", true, "Whether to enable transparent response compression")
	useTls                = flag.Bool("tls", false, "Whether to enable TLS")
	tlsCert               = flag.String("cert", "", "Full certificate file path")
	tlsKey                = flag.String("key", "", "Full key file path")
	unknown403            = flag.Bool("unknown403", true, "Return 403 on unknown paths.")
	authToken             = []byte(readFileUnsafe("token"))
	pathRoot              = []byte("/")
	pathFavicon           = []byte("/favicon.ico")
	htmlFile              = readFileUnsafe("www/index.html")
	lastBeat, missingBeat = readLastBeatSafe()
	gitCommitHash         = "A Development Version" // This is changed with compile flags in Makefile
)

func main() {
	flag.Parse()

	h := requestHandler
	if *compress {
		h = fasthttp.CompressHandler(h)
	}

	if *useTls && len(*tlsCert) > 0 && len(*tlsKey) > 0 {
		if err := fasthttp.ListenAndServeTLS(*addr, *tlsCert, *tlsKey, h); err != nil {
			log.Fatalf("- Error in ListenAndServe: %s", err)
		}
	} else {
		if err := fasthttp.ListenAndServe(*addr, h); err != nil {
			log.Fatalf("- Error in ListenAndServe: %s", err)
		}
	}
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set(fasthttp.HeaderServer, "Living's Heartbeat")

	path := ctx.Path()
	requestPath := string(path)

	if bytes.Equal(path, pathRoot) {
		handleRoot(ctx)
		return
	}

	if !ctx.IsGet() { // Don't allow get on non-root
		handleUnknown(ctx)
		return
	}

	if bytes.Equal(path, pathFavicon) {
		handleFavicon(ctx)
		return
	}

	pathInFolder := []string{"www", requestPath}
	pathFile := strings.Join(pathInFolder, "") // sep is not / because requestPath is prefixed with a /

	err := serveFile(ctx, pathFile, false)

	// Try to find corresponding file with the .html suffix added
	if err != nil {
		pathInFolder = []string{pathFile, ".html"}
		pathFile = strings.Join(pathInFolder, "")

		_ = serveFile(ctx, pathFile, true)
	}
}

func serveFile(ctx *fasthttp.RequestCtx, file string, handleErr bool) error {
	content, err := readFile(file)

	if err == nil {
		switch {
		case strings.HasSuffix(file, ".css"):
			ctx.Response.Header.Set(fasthttp.HeaderContentType, "text/css; charset=utf-8")
		case strings.HasSuffix(file, ".html"):
			ctx.Response.Header.Set(fasthttp.HeaderContentType, "text/html; charset=utf-8")
		}

		fmt.Fprint(ctx, content)
	} else if handleErr {
		handleUnknown(ctx)
	}

	return err
}
func handleRoot(ctx *fasthttp.RequestCtx) {
	// Serve the HTML page if it is a GET request
	if ctx.IsGet() {
		ctx.Response.Header.Set(fasthttp.HeaderContentType, "text/html; charset=utf-8")
		fmt.Fprint(ctx, getHtml())
		return
	}

	// We only want to allow POST requests if we're not serving a page
	if !ctx.IsPost() {
		ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
		fmt.Fprint(ctx, "400 Bad Request\n")
		log.Printf("- Returned 400 to %s - tried to connect with '%s'", ctx.RemoteIP(), ctx.Method())
		return
	}

	// The authentication key provided with said Auth header
	header := ctx.Request.Header.Peek("Auth")

	// Make sure Auth key is correct
	if !bytes.Equal(header, authToken) {
		handleUnknown(ctx)
		return
	}

	// Now that we know it is a POST and the Auth key is correct, return the unix timestamp and write it to a file
	handleBeat(ctx)
}

func handleUnknown(ctx *fasthttp.RequestCtx) {
	if *unknown403 {
		ctx.Response.SetStatusCode(fasthttp.StatusForbidden)
		fmt.Fprint(ctx, "403 Forbidden\n")
		log.Printf("- Returned 403 to %s - tried to connect with '%s' to '%s'", ctx.RemoteIP(), ctx.Request.Header.Peek("Auth"), ctx.Path())
	}
}

func handleFavicon(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set(fasthttp.HeaderContentType, "image/x-icon")
	f, _ := os.Open("www/favicon.ico")
	defer f.Close()
	io.Copy(ctx, f)
}

// Dynamic html is annoying so just replace a dummy value
func getHtml() string {
	lastBeatFormatted := intToTime(lastBeat).Format(time.RFC822)
	lastBeatStr := strconv.FormatInt(lastBeat, 10)

	currentTime := time.Now()
	currentBeatDifference := currentTime.Unix() - lastBeat

	// We also want to live update the current difference, instead of only when receiving a beat.
	if currentBeatDifference > missingBeat {
		missingBeat = currentBeatDifference
	}

	lastBeatArr := []string{lastBeatStr, lastBeatFormatted}
	timeDifference := timeDifference(lastBeat, currentTime)
	formattedAbsence := formattedTime(int(missingBeat))

	htmlTemp := strings.Replace(htmlFile, "LAST_BEAT", strings.Join(lastBeatArr, " "), 1)
	htmlTemp = strings.Replace(htmlTemp, "RELATIVE_TIME", timeDifference, 1)
	htmlTemp = strings.Replace(htmlTemp, "LONGEST_ABSENCE", formattedAbsence, 1)
	htmlTemp = strings.Replace(htmlTemp, "GIT_HASH", gitCommitHash, 2)
	htmlTemp = strings.Replace(htmlTemp, "GIT_REPO", "https://github.com/l1ving/heartbeat", 2)

	return htmlTemp
}

func handleBeat(ctx *fasthttp.RequestCtx) {
	newLastBeat := time.Now().Unix()
	currentBeatDifference := newLastBeat - lastBeat

	if currentBeatDifference > missingBeat {
		missingBeat = currentBeatDifference
	}

	lastBeatStr := strconv.FormatInt(newLastBeat, 10)
	missingBeatStr := strconv.FormatInt(missingBeat, 10)
	formattedBeats := []string{lastBeatStr, missingBeatStr}

	fmt.Fprintf(ctx, "%v\n", lastBeatStr)
	log.Printf("- Successful beat from %s", ctx.RemoteIP())

	lastBeat = newLastBeat
	writeToFile("last_beat", strings.Join(formattedBeats, ":"))
}
