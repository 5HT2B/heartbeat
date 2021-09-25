//go:generate qtc -dir=templates
package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/Ferluci/fast-realip"
	"github.com/technically-functional/heartbeat/templates"
	"github.com/valyala/fasthttp"
	"log"
	"strconv"
	"time"
)

var (
	addr                              = flag.String("addr", "localhost:6060", "TCP address to listen to")
	serverName                        = flag.String("name", "Liv's Heartbeat", "The name of the server to use")
	authToken                         = []byte(ReadFileUnsafe("config/token"))
	apiPrefix                         = []byte("/api/")
	cssPrefix                         = []byte("/css/")
	icoSuffix                         = []byte(".ico")
	pngSuffix                         = []byte(".png")
	gitCommitHash                     = "A Development Version" // This is changed with compile flags in Makefile
	timeFormat                        = "Jan 02 2006 15:04:05 MST"
	lastBeat, missingBeat, totalBeats = ReadLastBeatSafe()
	totalVisits                       = ReadGetRequestsSafe()
	gitRepo                           = "https://github.com/technically-functional/heartbeat"
	imgHandler                        = fasthttp.FSHandler("www", 0)
	cssHandler                        = fasthttp.FSHandler("www", 1)
)

//goland:noinspection HttpUrlsUsage
func main() {
	flag.Parse()
	log.Printf("- Running heartbeat on http://" + *addr)

	h := requestHandler
	h = fasthttp.CompressHandler(h)

	if err := fasthttp.ListenAndServe(*addr, h); err != nil {
		log.Fatalf("- Error in ListenAndServe: %s", err)
	}
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set(fasthttp.HeaderServer, *serverName)

	path := ctx.Path()
	pathStr := string(path)

	switch {
	case bytes.HasPrefix(path, apiPrefix):
		apiHandler(ctx, pathStr)
	case bytes.HasPrefix(path, cssPrefix):
		cssHandler(ctx)
	case bytes.HasSuffix(path, icoSuffix), bytes.HasSuffix(path, pngSuffix):
		imgHandler(ctx)
	default:
		totalVisits++
		ctx.Response.Header.Set(fasthttp.HeaderContentType, "text/html; charset=utf-8")

		switch pathStr {
		case "/":
			mainPageHandler(ctx)
		case "/privacy":
			privacyPolicyPageHandler(ctx)
		case "/stats":
			statsPageHandler(ctx)
		default:
			errorPageHandler(ctx, fasthttp.StatusNotFound, "404 Not Found")
		}
	}
}

func apiHandler(ctx *fasthttp.RequestCtx, path string) {
	if !ctx.IsPost() {
		errorPageHandler(ctx, fasthttp.StatusBadRequest, "400 Bad Request")
		return
	}

	// The authentication key provided with said Auth header
	header := ctx.Request.Header.Peek("Auth")

	// Make sure Auth key is correct
	if !bytes.Equal(header, authToken) {
		errorPageHandler(ctx, fasthttp.StatusForbidden, "403 Forbidden")
		return
	}

	switch path {
	case "/api/beat":
		handleSuccessfulBeat(ctx)
	default:
		errorPageHandler(ctx, fasthttp.StatusBadRequest, "400 Bad Request")
	}
}

func mainPageHandler(ctx *fasthttp.RequestCtx) {
	p := getMainPage()
	templates.WritePageTemplate(ctx, p)
}
func privacyPolicyPageHandler(ctx *fasthttp.RequestCtx) {
	p := &templates.PrivacyPolicyPage{
		ServerName: *serverName,
	}
	templates.WritePageTemplate(ctx, p)
}

func statsPageHandler(ctx *fasthttp.RequestCtx) {
	p := &templates.StatsPage{
		TotalBeats:   FormattedNum(totalBeats),
		TotalDevices: FormattedNum(2), // TODO: Add support for this
		TotalVisits:  FormattedNum(totalVisits),
		ServerName:   *serverName,
	}
	templates.WritePageTemplate(ctx, p)
}

func errorPageHandler(ctx *fasthttp.RequestCtx, code int, message string) {
	p := &templates.ErrorPage{
		Message: message,
		Path:    ctx.Path(),
		Method:  ctx.Method(),
	}
	templates.WritePageTemplate(ctx, p)
	ctx.SetStatusCode(code)
	log.Printf("- Returned %v to %s - tried to connect with '%s'", code, realip.FromRequest(ctx), ctx.Method())
}

func getMainPage() *templates.MainPage {
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

	page := &templates.MainPage{
		LastSeen:       lastSeen,
		TimeDifference: timeDifference,
		MissingBeat:    missingBeatFmt,
		TotalBeats:     totalBeatsFmt,
		CurrentTime:    currentTimeStr,
		GitHash:        gitCommitHash,
		GitRepo:        gitRepo,
		ServerName:     *serverName,
	}

	return page
}

func handleSuccessfulBeat(ctx *fasthttp.RequestCtx) {
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
