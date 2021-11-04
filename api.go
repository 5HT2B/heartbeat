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
	apiBeatPath          = "/api/beat"
	apiUpdateStatsPath   = "/api/update/stats"
	apiUpdateDevicesPath = "/api/update/devices"
	apiInfoPath          = "/api/info"
	apiStatsPath         = "/api/stats"
	apiDevicesPath       = "/api/devices"
	jsonMime             = "application/json"
)

func ApiHandler(ctx *fasthttp.RequestCtx, path string) {

	switch path {
	case apiBeatPath, apiUpdateStatsPath, apiUpdateDevicesPath:
		if !ctx.IsPost() {
			ErrorBadRequest(ctx, true)
			return
		}

		// The authentication key provided with said Auth header
		header := ctx.Request.Header.Peek("Auth")

		// Make sure Auth key is correct
		if string(header) != authToken {
			ErrorForbidden(ctx, true)
			return
		}
	case apiInfoPath, apiStatsPath, apiDevicesPath:
		heartbeatStats.TotalVisits += 1
		if !ctx.IsGet() {
			ErrorBadRequest(ctx, true)
			return
		}
	}

	switch path {
	case apiBeatPath:
		HandleSuccessfulBeat(ctx)
	case apiUpdateStatsPath:
		handleUpdateStats(ctx)
	case apiUpdateDevicesPath:
		handleUpdateDevices(ctx)
	case apiInfoPath:
		handleJsonObject(ctx, FormattedInfo())
	case apiStatsPath:
		UpdateUptime()
		UpdateLastBeatFmt()
		handleJsonObject(ctx, heartbeatStats)
	case apiDevicesPath:
		handleJsonObject(ctx, heartbeatDevices)
	default:
		ErrorBadRequest(ctx, true)
	}
}

// handleUpdateStats will allow authenticated users to replace the current stats
func handleUpdateStats(ctx *fasthttp.RequestCtx) {
	var stats HeartbeatStats
	err := json2.Unmarshal(ctx.PostBody(), &stats)
	if err != nil {
		HandleClientErr(ctx, "Error unmarshalling json", err)
		return
	}

	heartbeatStats = &stats
	HandleSuccess(ctx)
}

// handleUpdateDevices will allow authenticated users to replace the current devices
func handleUpdateDevices(ctx *fasthttp.RequestCtx) {
	var devices []HeartbeatDevice
	err := json2.Unmarshal(ctx.PostBody(), &devices)
	if err != nil {
		HandleClientErr(ctx, "Error unmarshalling json", err)
		return
	}

	heartbeatDevices = &devices
	HandleSuccess(ctx)
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

// HandleSuccessfulBeat will update the db's Device with the new last beat, and print the last beat's unix timestamp
func HandleSuccessfulBeat(ctx *fasthttp.RequestCtx) {
	device := ctx.Request.Header.Peek("Device")
	// Make sure a device is set
	if len(device) == 0 {
		ErrorBadRequest(ctx, true)
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
