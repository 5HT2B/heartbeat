package main

import (
	json2 "encoding/json"
	"fmt"
	"github.com/Ferluci/fast-realip"
	"github.com/valyala/fasthttp"
	"log"
	"strconv"
	"time"
)

var (
	apiBeatPath    = "/api/beat"
	apiStatsPath   = "/api/stats"
	apiDevicesPath = "/api/devices"
	jsonMime       = "application/json"
)

func ApiHandler(ctx *fasthttp.RequestCtx, path string) {
	switch path {
	case apiBeatPath:
		if !ctx.IsPost() {
			ErrorBadRequest(ctx)
			return
		}

		// The authentication key provided with said Auth header
		header := ctx.Request.Header.Peek("Auth")

		// Make sure Auth key is correct
		if string(header) != authToken {
			ErrorForbidden(ctx)
			return
		}
	case apiStatsPath, apiDevicesPath:
		if !ctx.IsGet() {
			ErrorBadRequest(ctx)
			return
		}
	}

	switch path {
	case apiBeatPath:
		HandleSuccessfulBeat(ctx)
	case apiStatsPath:
		UpdateUptime()
		handleJsonObject(ctx, heartbeatStats)
	case apiDevicesPath:
		handleJsonObject(ctx, heartbeatDevices)
	default:
		ErrorBadRequest(ctx)
	}
}

// handleJsonObject will print the raw json of v
func handleJsonObject(ctx *fasthttp.RequestCtx, v interface{}) {
	json, err := json2.MarshalIndent(v, "", "    ")

	if err != nil {
		HandleInternalErr(ctx, "Error formatting json", err)
		return
	}

	ctx.Response.Header.SetContentType(jsonMime)
	_, _ = fmt.Fprintf(ctx, "%s\n", json)
}

func HandleSuccessfulBeat(ctx *fasthttp.RequestCtx) {
	device := ctx.Request.Header.Peek("Device")
	// Make sure a device is set
	if len(device) == 0 {
		ErrorBadRequest(ctx)
		return
	}

	timestamp := time.Now().Unix()
	timestampStr := strconv.FormatInt(timestamp, 10)

	err := UpdateLastBeat(string(device), timestamp)
	if err != nil {
		HandleInternalErr(ctx, "Error updating last beat", err)
		return
	}

	_, _ = fmt.Fprintf(ctx, "%v\n", timestampStr)
	log.Printf("- Successful beat from %s", realip.FromRequest(ctx))
}
