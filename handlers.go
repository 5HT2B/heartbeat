package main

import (
	"bytes"
	"github.com/5HT2B/heartbeat/templates"
	"github.com/valyala/fasthttp"
	"strings"
	"time"
)

var (
	apiPrefix = []byte("/api/")
	cssPrefix = []byte("/css/")
	icoSuffix = []byte(".ico")
	pngSuffix = []byte(".png")
	gitRepo   = "https://github.com/5HT2B/heartbeat" // set in .env

	cssHandler = fasthttp.FSHandler("www", 1)
	imgHandler = fasthttp.FSHandler("www", 0)
)

func RequestHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set(fasthttp.HeaderServer, serverName)

	path := ctx.Path()
	pathStr := string(path)
	if pathStr != "/" {
		pathStr = strings.TrimSuffix(pathStr, "/")
	}

	switch {
	case bytes.HasPrefix(path, apiPrefix):
		ApiHandler(ctx, pathStr)
	case bytes.HasPrefix(path, cssPrefix):
		cssHandler(ctx)
	case bytes.HasSuffix(path, icoSuffix), bytes.HasSuffix(path, pngSuffix):
		imgHandler(ctx)
	default:
		heartbeatStats.TotalVisits += 1
		heartbeatStats.TotalVisitsFormatted = FormattedNum(heartbeatStats.TotalVisits)
		ctx.Response.Header.Set(fasthttp.HeaderContentType, "text/html; charset=utf-8")

		switch pathStr {
		case "/":
			MainPageHandler(ctx)
		case "/privacy":
			PrivacyPolicyPageHandler(ctx)
		case "/stats":
			StatsPageHandler(ctx)
		default:
			ErrorNotFound(ctx, false)
		}
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
	UpdateNumDevices()
	p := &templates.StatsPage{
		TotalBeats:   heartbeatStats.TotalBeatsFormatted,
		TotalDevices: heartbeatStats.TotalDevicesFormatted,
		TotalVisits:  heartbeatStats.TotalVisitsFormatted,
		TotalUptime:  heartbeatStats.TotalUptimeFormatted,
		ServerName:   serverName,
	}
	templates.WritePageTemplate(ctx, p)
}

func getMainPage() *templates.MainPage {
	currentTime := time.Now()
	lastBeat := GetLastBeat()
	UpdateLastBeatFmtV(lastBeat, currentTime)

	page := &templates.MainPage{
		LastSeen:       LastSeen(),                       // date last seen
		TimeDifference: heartbeatStats.LastBeatFormatted, // relative time to last seen
		MissingBeat:    LongestAbsence(),                 // longest absence
		TotalBeats:     heartbeatStats.TotalBeatsFormatted,
		CurrentTime:    currentTime.Format(timeFormat),
		GitHash:        gitCommitHash,
		GitRepo:        gitRepo,
		ServerName:     serverName,
	}

	return page
}
