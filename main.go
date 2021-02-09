package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/valyala/fasthttp"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	addr        = flag.String("addr", ":6060", "TCP address to listen to")
	compress    = flag.Bool("compress", true, "Whether to enable transparent response compression")
	authToken   = []byte(readFileUnsafe("token"))
	pathRoot    = []byte("/")
	pathStyle   = []byte("/style.css")
	pathFavicon = []byte("/favicon.ico")
	htmlFile    = readFileUnsafe("index.html")
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

	base10Time := strconv.FormatInt(lastBeat, 10)

	fmt.Fprintf(ctx, "%v\n", base10Time)
	log.Printf("- Successful beat from %s", ctx.RemoteIP())

	writeToFile("last_beat", base10Time)
	writeToFile("last_beat_formatted", lastTime.Format(time.RFC822))
}

func handleUnknown(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetStatusCode(fasthttp.StatusForbidden)
	fmt.Fprint(ctx, "403 Forbidden\n")
	log.Printf("- Returned 403 to %s - tried to connect with '%s' to '%s'", ctx.RemoteIP(), ctx.Request.Header.Peek("Auth"), ctx.Path())
}

func handleStyle(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set(fasthttp.HeaderContentType, "text/css; charset=utf-8")
	fmt.Fprint(ctx, readFileUnsafe("style.css"))
}

func handleFavicon(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set(fasthttp.HeaderContentType, "image/x-icon")
	f, _ := os.Open("favicon.ico")
	defer f.Close()
	io.Copy(ctx, f)
}

// Dynamic html is annoying so just replace a dummy value lol
func getHtml() string {
	lastBeat, err1 := readFile("last_beat")
	lastBeatFormatted, err2 := readFile("last_beat_formatted")

	if err1 != nil {
		lastBeat = "Error reading last_beat from server, this should not happen."
	}

	if err2 != nil {
		lastBeatFormatted = "Error reading last_beat_formatted from server, this should not happen."
	}

	formattedTime := []string{lastBeat, lastBeatFormatted}

	htmlWithBeat := strings.Replace(htmlFile, "LAST_BEAT", strings.Join(formattedTime, " "), 1)
	timeDifference := timeDifference(lastBeat, time.Now())

	return strings.Replace(htmlWithBeat, "RELATIVE_TIME", timeDifference, 1)
}

func readFileUnsafe(file string) string {
	content, err := readFile(file)

	if err != nil {
		log.Printf("- Failed to read '%s'", file)
	}

	return content
}

func readFile(file string) (string, error) {
	dat, err := ioutil.ReadFile(file)
	return string(dat), err
}

func writeToFile(file string, content string) {
	data := []byte(content)
	err := ioutil.WriteFile(file, data, 0644)

	if err != nil {
		log.Printf("- Failed to read '%s'", file)
	}
}

func timeDifference(lastBeat string, t time.Time) string {
	lastBeatInt, err := strconv.ParseInt(lastBeat, 10, 64)

	if err != nil {
		return fmt.Sprintf("Could not convert '%s' to int64", lastBeat)
	}

	newTime := time.Unix(lastBeatInt, 0)
	st := t.Sub(newTime)

	return formattedTime(int(st.Seconds()))
}

func formattedTime(secondsIn int) string {
	hours := secondsIn / 3600
	minutes := (secondsIn / 60) - (60 * hours)
	seconds := secondsIn % 60
	str := fmt.Sprintf("%d hours, %d minutes, %d seconds", hours, minutes, seconds)
	return str
}
