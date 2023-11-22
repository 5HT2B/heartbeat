package main

import (
	json2 "encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/fasthttp/websocket"
	realip "github.com/ferluci/fast-realip"
	"github.com/valyala/fasthttp"
)

var (
	apiBeatPath          = "/api/beat"
	apiUpdateStatsPath   = "/api/update/stats"
	apiUpdateDevicesPath = "/api/update/devices"
	apiInfoPath          = "/api/info"
	apiStatsPath         = "/api/stats"
	apiDevicesPath       = "/api/devices"
	apiRealtimeStatsPath = "/api/realtime/stats"
	jsonMime             = "application/json"
)

var upgrader = websocket.FastHTTPUpgrader{}

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
		if !ctx.IsGet() {
			ErrorBadRequest(ctx, true)
			return
		}
		if string(ctx.QueryArgs().Peek("countVisit")) != "false" {
			heartbeatStats.TotalVisits += 1
			heartbeatStats.TotalVisitsFormatted = FormattedNum(heartbeatStats.TotalVisits)
		}
	case apiRealtimeStatsPath:
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
		UpdateNumDevices()
		UpdateLastBeatFmt()
		handleJsonObject(ctx, heartbeatStats)
	case apiRealtimeStatsPath:
		handleWs(ctx, heartbeatStats)
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

func handleWs(ctx *fasthttp.RequestCtx, stats *HeartbeatStats) {
	err := upgrader.Upgrade(ctx, func(conn *websocket.Conn) {
		defer conn.Close()
		for {
			UpdateUptime()
			UpdateNumDevices()
			UpdateLastBeatFmt()
			json, err := json2.Marshal(stats)
			if err != nil {
				HandleInternalErr(ctx, "Error formatting json", err)
				return
			}
			err = conn.WriteMessage(websocket.TextMessage, json)
			if err != nil {
				log.Printf("writing to websocket: %v", err)
				break
			}
			time.Sleep(time.Second)
		}
	})
	if err != nil {
		// FIXME: handle this
		return
	}
}
