package main

import (
	"bytes"
	"fmt"
	"github.com/Ferluci/fast-realip"
	"github.com/technically-functional/heartbeat/templates"
	"github.com/valyala/fasthttp"
	"log"
	"strconv"
	"time"
)

var (
	apiPrefix = []byte("/api/")
	cssPrefix = []byte("/css/")
	icoSuffix = []byte(".ico")
	pngSuffix = []byte(".png")
	gitRepo   = "https://github.com/technically-functional/heartbeat" // set in .env

	cssHandler = fasthttp.FSHandler("www", 1)
	imgHandler = fasthttp.FSHandler("www", 0)
)

func RequestHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set(fasthttp.HeaderServer, serverName)

	path := ctx.Path()
	pathStr := string(path)

	switch {
	case bytes.HasPrefix(path, apiPrefix):
		ApiHandler(ctx, pathStr)
	case bytes.HasPrefix(path, cssPrefix):
		cssHandler(ctx)
	case bytes.HasSuffix(path, icoSuffix), bytes.HasSuffix(path, pngSuffix):
		imgHandler(ctx)
	default:
		heartbeatStats.TotalVisits += 1
		ctx.Response.Header.Set(fasthttp.HeaderContentType, "text/html; charset=utf-8")

		switch pathStr {
		case "/":
			MainPageHandler(ctx)
		case "/privacy":
			PrivacyPolicyPageHandler(ctx)
		case "/stats":
			StatsPageHandler(ctx)
		default:
			ErrorPageHandler(ctx, fasthttp.StatusNotFound, "404 Not Found", false)
		}
	}
}

func ApiHandler(ctx *fasthttp.RequestCtx, path string) {
	if !ctx.IsPost() {
		ErrorPageHandler(ctx, fasthttp.StatusBadRequest, "400 Bad Request", true)
		return
	}

	// The authentication key provided with said Auth header
	header := ctx.Request.Header.Peek("Auth")
	device := ctx.Request.Header.Peek("Device")

	// Make sure Auth key is correct
	if string(header) != authToken {
		ErrorPageHandler(ctx, fasthttp.StatusForbidden, "403 Forbidden", true)
		return
	}

	// Make sure a device is set
	if len(device) == 0 {
		ErrorPageHandler(ctx, fasthttp.StatusBadRequest, "400 Bad Request", true)
		return
	}

	switch path {
	case "/api/beat":
		HandleSuccessfulBeat(ctx, string(device))
	default:
		ErrorPageHandler(ctx, fasthttp.StatusBadRequest, "400 Bad Request", true)
	}
}

func MainPageHandler(ctx *fasthttp.RequestCtx) {
	p := getMainPage()
	templates.WritePageTemplate(ctx, p)
}

func PrivacyPolicyPageHandler(ctx *fasthttp.RequestCtx) {
	p := &templates.PrivacyPolicyPage{
		ServerName: serverName,
	}
	templates.WritePageTemplate(ctx, p)
}

func StatsPageHandler(ctx *fasthttp.RequestCtx) {
	UpdateUptime()
	p := &templates.StatsPage{
		TotalBeats:   FormattedNum(heartbeatStats.TotalBeats),
		TotalDevices: FormattedNum(int64(len(*heartbeatDevices))),
		TotalVisits:  FormattedNum(heartbeatStats.TotalVisits),
		TotalUptime:  FormattedTime(heartbeatStats.TotalUptime),
		ServerName:   serverName,
	}
	templates.WritePageTemplate(ctx, p)
}

func ErrorPageHandler(ctx *fasthttp.RequestCtx, code int, message string, plaintext bool) {
	ctx.SetStatusCode(code)
	log.Printf("- Returned %v to %s - tried to connect with %s to %s",
		code, realip.FromRequest(ctx), ctx.Method(), ctx.Path())

	if plaintext {
		ctx.Response.Header.Set(fasthttp.HeaderContentType, "text/plain; charset=utf-8")
		_, _ = fmt.Fprintf(ctx, "%v %s\n", code, message)
	} else {
		p := &templates.ErrorPage{
			Message: message,
			Path:    ctx.Path(),
			Method:  ctx.Method(),
		}
		templates.WritePageTemplate(ctx, p)
	}
}

func getMainPage() *templates.MainPage {
	currentTime := time.Now()
	lastBeat := GetLastBeat()
	timeDifference := "Never"
	if lastBeat != nil {
		timeDifference = TimeDifference(lastBeat.Timestamp, currentTime)
	}

	page := &templates.MainPage{
		LastSeen:       LastSeen(),       // date last seen
		TimeDifference: timeDifference,   // relative time to last seen
		MissingBeat:    LongestAbsence(), // longest absence
		TotalBeats:     FormattedNum(heartbeatStats.TotalBeats),
		CurrentTime:    currentTime.Format(timeFormat),
		GitHash:        gitCommitHash,
		GitRepo:        gitRepo,
		ServerName:     serverName,
	}

	return page
}

func HandleSuccessfulBeat(ctx *fasthttp.RequestCtx, device string) {
	timestamp := time.Now().Unix()
	timestampStr := strconv.FormatInt(timestamp, 10)

	err := UpdateLastBeat(device, timestamp)
	if err != nil {
		ErrorPageHandler(ctx, fasthttp.StatusInternalServerError, err.Error(), true)
		return
	}

	_, _ = fmt.Fprintf(ctx, "%v\n", timestampStr)
	log.Printf("- Successful beat from %s", realip.FromRequest(ctx))
}
