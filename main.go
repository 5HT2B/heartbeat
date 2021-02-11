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
	addr        = flag.String("addr", ":6060", "TCP address to listen to")
	compress    = flag.Bool("compress", true, "Whether to enable transparent response compression")
	useTls      = flag.Bool("tls", false, "Whether to enable TLS")
	tlsCert     = flag.String("cert", "", "Full certificate file path")
	tlsKey      = flag.String("key", "", "Full key file path")
	unknown403  = flag.Bool("unknown403", true, "Return 403 on unknown paths.")
	authToken   = []byte(readFileUnsafe("token"))
	pathRoot    = []byte("/")
	pathStyle   = []byte("/style.css")
	pathFavicon = []byte("/favicon.ico")
	htmlFile    = readFileUnsafe("index.html")
	cssFile     = readFileUnsafe("style.css")
	lastBeat    = readLastBeatSafe()
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
	switch {
	case bytes.Equal(path, pathRoot):
		handleRoot(ctx)
	case bytes.Equal(path, pathStyle):
		handleStyle(ctx)
	case bytes.Equal(path, pathFavicon):
		handleFavicon(ctx)
	default:
		handleUnknown(ctx)
	}

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
	newLastBeat := time.Now().Unix()
	base10Time := strconv.FormatInt(newLastBeat, 10)

	fmt.Fprintf(ctx, "%v\n", base10Time)
	log.Printf("- Successful beat from %s", ctx.RemoteIP())

	lastBeat = newLastBeat
	writeToFile("last_beat", base10Time)
}

func handleUnknown(ctx *fasthttp.RequestCtx) {
	if *unknown403 {
		ctx.Response.SetStatusCode(fasthttp.StatusForbidden)
		fmt.Fprint(ctx, "403 Forbidden\n")
		log.Printf("- Returned 403 to %s - tried to connect with '%s' to '%s'", ctx.RemoteIP(), ctx.Request.Header.Peek("Auth"), ctx.Path())
	}
}

func handleStyle(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set(fasthttp.HeaderContentType, "text/css; charset=utf-8")
	fmt.Fprint(ctx, cssFile)
}

func handleFavicon(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set(fasthttp.HeaderContentType, "image/x-icon")
	f, _ := os.Open("favicon.ico")
	defer f.Close()
	io.Copy(ctx, f)
}

// Dynamic html is annoying so just replace a dummy value lol
func getHtml() string {
	lastBeatFormatted := intToTime(lastBeat).Format(time.RFC822)
	lastBeatStr := strconv.FormatInt(lastBeat, 10)

	formattedTime := []string{lastBeatStr, lastBeatFormatted}

	htmlWithBeat := strings.Replace(htmlFile, "LAST_BEAT", strings.Join(formattedTime, " "), 1)
	timeDifference := timeDifference(lastBeatStr, time.Now())

	return strings.Replace(htmlWithBeat, "RELATIVE_TIME", timeDifference, 1)
}
