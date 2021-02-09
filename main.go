package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	addr        = flag.String("addr", ":6060", "TCP address to listen to")
	compress    = flag.Bool("compress", true, "Whether to enable transparent response compression")
	authToken   = []byte(readFile("token"))
	pathRoot    = []byte("/")
	pathStyle   = []byte("/style.css")
	pathFavicon = []byte("/favicon.ico")
	htmlFile    = readFile("index.html")
)

func main() {
	flag.Parse()

	h := requestHandler
	if *compress {
		h = fasthttp.CompressHandler(h)
	}

	if err := fasthttp.ListenAndServe(*addr, h); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
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
		ctx.Response.SetStatusCode(fasthttp.StatusForbidden)
		fmt.Fprint(ctx, "403 Forbidden\n")
		log.Printf("- Returned 403 to %s - tried to connect with '%s'", ctx.RemoteIP(), header)
		return
	}

	// Now that we know it is a POST and the Auth key is correct, return the unix timestamp and write it to a file
	lastTime := time.Now()
	lastBeat := lastTime.Unix()

	formattedTime := []string{strconv.FormatInt(lastBeat, 10), lastTime.Format(time.RFC822)}

	fmt.Fprintf(ctx, "%v\n", formattedTime)
	log.Printf("- Successful beat from %s", ctx.RemoteIP())
	writeToFile("last_beat", strings.Join(formattedTime, " "))
}

func handleUnknown(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetStatusCode(fasthttp.StatusForbidden)
	fmt.Fprint(ctx, "403 Forbidden\n")
	log.Printf("- Returned 403 to %s - tried to connect with '%s' to '%s'", ctx.RemoteIP(), ctx.Request.Header.Peek("Auth"), ctx.Path())
}

func handleStyle(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set(fasthttp.HeaderContentType, "text/css; charset=utf-8")
	fmt.Fprint(ctx, readFile("style.css"))
}

// TODO: Figure out how files work via htp
func handleFavicon(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set(fasthttp.HeaderContentType, "text/html; charset=utf-8")

}

// Dynamic html is annoying so just replace a dummy value lol
func getHtml() string {
	lastBeat := readFile("last_beat")
	return strings.Replace(htmlFile, "LAST_BEAT", lastBeat, 1)
}

// Read file and log if it errored
func readFile(file string) string {
	dat, err := ioutil.ReadFile(file)

	if err != nil {
		log.Printf("- Failed to read '%s'", file)
	}

	return string(dat)
}

func writeToFile(file string, content string) {
	data := []byte(content)
	err := ioutil.WriteFile(file, data, 0644)

	if err != nil {
		log.Printf("- Failed to read '%s'", file)
	}
}
