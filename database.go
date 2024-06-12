package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/nitishm/go-rejson/v4"
)

var (
	rdb *redis.Client   // set in main()
	rjh *rejson.Handler // set in main()

	redisAddr = "database:6379" // set in .env
	redisPass = ""              // set in .env

	heartbeatStats   *HeartbeatStats    = nil
	heartbeatDevices *[]HeartbeatDevice = nil
)

// SetupDatabase creates the ReJSON handler and Redis client
func SetupDatabase() (*redis.Client, *rejson.Handler) {
	rh := rejson.NewReJSONHandler()
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPass,
	})

	// If connection doesn't work, panic
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("- Failed to ping Redis server: %v\n", err)
	}

	// We have a working connection
	log.Printf("- Connected to Redis at %s", redisAddr)

	rh.SetGoRedisClient(client)
	return client, rh
}

// SetupDatabaseSaving will run SaveLocalInDatabase every minute with a ticket
func SetupDatabaseSaving() {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			SaveLocalInDatabase()
		}
	}()
}

// SetupLocalValues will look for existing stats and devices in the database, and read them to local values
func SetupLocalValues() {
	if res, err := rjh.JSONGet("stats", "."); err != nil {
		log.Printf("- Error loading stats (using default): %v\n", err)
		heartbeatStats = defaultHeartbeatStats
	} else {
		var stats HeartbeatStats
		if err = json.Unmarshal(res.([]byte), &stats); err != nil {
			panic(err)
		}
		log.Printf("- Loaded stats: %v\n", stats)
		heartbeatStats = &stats
	}

	if res, err := rjh.JSONGet("devices", "."); err != nil {
		log.Printf("- Error loading devices (using default): %v\n", err)
		heartbeatDevices = defaultHeartbeatDevices
	} else {
		var devices []HeartbeatDevice
		if err = json.Unmarshal(res.([]byte), &devices); err != nil {
			panic(err)
		}
		log.Printf("- Loaded devices: %v\n", devices)
		heartbeatDevices = &devices
	}
}

// SaveLocalInDatabase will save local copies of small database stats and devices to the database every 5 minutes
func SaveLocalInDatabase() {
	// This is also called when viewing the stats page, but we want to run it here to avoid missing time
	// if nobody looks at the stats page
	UpdateUptime()

	log.Printf("- Autosaving database, uptime is %v\n", heartbeatStats.TotalUptime)

	if _, err := rjh.JSONSet("stats", ".", heartbeatStats); err != nil {
		log.Fatalf("Error saving stats: %v", err)
	}

	if _, err := rjh.JSONSet("devices", ".", heartbeatDevices); err != nil {
		log.Fatalf("Error saving devices: %v", err)
	}
}

// FormattedInfo will return the formatted info, displayed on the main page
func FormattedInfo() HeartbeatInfo {
	currentTime := time.Now()
	lastBeat := GetLastBeat()
	UpdateLastBeatFmtV(lastBeat, currentTime)
	return HeartbeatInfo{
		LastSeen:       LastSeen(),
		TimeDifference: heartbeatStats.LastBeatFormatted,
		MissingBeat:    LongestAbsence(),
		TotalBeats:     heartbeatStats.TotalBeatsFormatted,
	}
}

// LastSeen will return the formatted date of the last timestamp received from a beat
func LastSeen() string {
	lastBeat := GetLastBeat()
	if lastBeat == nil {
		return "Never"
	}
	return time.Unix(lastBeat.Timestamp, 0).Format(timeFormat)
}

// LongestAbsence will return HeartbeatStats.LongestMissingBeat unless the current missing beat is longer
func LongestAbsence() string {
	lastBeat := GetMostRecentBeat()

	// This will happen when GetLastBeat returned a nil, and heartbeatDevices is empty
	if lastBeat == nil {
		heartbeatStats.LastBeatFormatted = "Never"
		return "Never"
	}

	diff := time.Now().Unix() - lastBeat.Timestamp
	// If current absence is bigger than record absence, return current absence
	if diff > heartbeatStats.LongestMissingBeat {
		heartbeatStats.LongestMissingBeat = diff
		return FormattedTime(diff)
	} else {
		return FormattedTime(heartbeatStats.LongestMissingBeat)
	}
}

func UpdateNumDevices() {
	heartbeatStats.TotalDevicesFormatted = FormattedNum(int64(len(*heartbeatDevices)))
}

// UpdateUptime will update the uptime statistics
func UpdateUptime() {
	now := time.Now().UnixMilli()
	diff := now - uptimeTimer
	uptimeTimer = now
	heartbeatStats.TotalUptime += diff
	heartbeatStats.TotalUptimeFormatted = FormattedTime(heartbeatStats.TotalUptime / 1000)
}

// UpdateLastBeatFmt will update the formatted last beat statistic
func UpdateLastBeatFmt() {
	currentTime := time.Now()
	lastBeat := GetLastBeat()
	if lastBeat != nil {
		heartbeatStats.LastBeatFormatted = TimeDifference(lastBeat.Timestamp, currentTime)
	}
}

// UpdateLastBeatFmtV will update the formatted last beat statistic
func UpdateLastBeatFmtV(lastBeat *HeartbeatBeat, currentTime time.Time) {
	if lastBeat != nil {
		heartbeatStats.LastBeatFormatted = TimeDifference(lastBeat.Timestamp, currentTime)
	}
}

// UpdateTotalBeats will update the formatted total beats statistic
func UpdateTotalBeats() {
	heartbeatStats.TotalBeats += 1
	heartbeatStats.TotalBeatsFormatted = FormattedNum(heartbeatStats.TotalBeats)
}

// UpdateDevice will update the LastBeat of a device
func UpdateDevice(beat HeartbeatBeat) {
	var device HeartbeatDevice
	n := -1
	for index, tmpDevice := range *heartbeatDevices {
		if tmpDevice.DeviceName == beat.DeviceName {
			device = tmpDevice
			n = index
			break // name should only ever be matching once
		}
	}

	if n == -1 { // couldn't find an existing device with the name
		device = HeartbeatDevice{beat.DeviceName, beat, 0, 0}
	}

	diff := time.Now().Unix() - device.LastBeat.Timestamp
	if diff > device.LongestMissingBeat {
		device.LongestMissingBeat = diff
	}
	// We want to update the longest absence here (heartbeatStats.LongestMissingBeat) in case
	// device.LongestMissingBeat > heartbeatStats.LongestMissingBeat *and* other devices haven't pinged recently
	LongestAbsence()

	device.LastBeat = beat
	device.TotalBeats += 1
	UpdateTotalBeats()

	if n == -1 { // if device doesn't exist, append it, else replace it
		*heartbeatDevices = append(*heartbeatDevices, device)
		PostMessage("New device added", fmt.Sprintf("A new device called `%s` was added on <t:%v:d> at <t:%v:T>", beat.DeviceName, beat.Timestamp, beat.Timestamp), EmbedColorGreen, WebhookLevelSimple)
	} else {
		(*heartbeatDevices)[n] = device
	}
}

// UpdateLastBeat will save the last beat and insert a new HeartbeatBeat into beats
func UpdateLastBeat(deviceName string, timestamp int64) error {
	oldLastBeat := GetMostRecentBeat()
	lastBeat := HeartbeatBeat{deviceName, timestamp}
	UpdateDevice(lastBeat)

	if _, err := rjh.JSONSet("last_beat", ".", lastBeat); err != nil {
		return err
	}

	lastBeatArr := []HeartbeatBeat{lastBeat}
	err := appendOrCreateArr("beats", ".", lastBeat, lastBeatArr)

	if err == nil {
		PostMessage("Successful beat", fmt.Sprintf("From `%s` on <t:%v:d> at <t:%v:T>", deviceName, timestamp, timestamp), EmbedColorBlue, WebhookLevelAll)

		if oldLastBeat != nil && time.Duration(timestamp-oldLastBeat.Timestamp)*time.Second > 1*time.Hour {
			PostMessage(
				"Absence longer than 1 hour",
				fmt.Sprintf(
					"From <t:%v> to <t:%v>\nUTC Data: `%s,%s`",
					oldLastBeat.Timestamp, timestamp,
					FormattedUTCData(oldLastBeat.Timestamp), FormattedUTCData(timestamp),
				),
				EmbedColorOrange, WebhookLevelLongAbsence,
			)
		}
	}
	return err
}

// GetMostRecentBeat will return the most recent beat regardless of device, not just last-inserted beat
func GetMostRecentBeat() *HeartbeatBeat {
	lastBeat := GetLastBeat()
	for _, device := range *heartbeatDevices {
		// This technically shouldn't be possible, as UpdateDevice is called inside UpdateLastBeat.
		// Nevertheless, we would rather avoid a random panic from accessing a nil reference.
		if lastBeat == nil {
			lastBeat = &device.LastBeat
		}
		// If this device has a more recent beat than the most recent beat's device, use it instead.
		// The reasoning behind this is, if we suddenly get a new beat from a device that hasn't sent a beat in a while
		// we don't want it to set the longest absence to the old device's oldest absence, as this new device
		// has sent beats more recently, and you have not actually been absent as long as the original lastBeat here.
		if device.LastBeat.Timestamp > lastBeat.Timestamp {
			lastBeat = &device.LastBeat
		}
	}

	return lastBeat
}

// GetLastBeat will get the last beat, and return nilLastBeat if there was an error retrieving it (likely no beat)
func GetLastBeat() *HeartbeatBeat {
	res, err := rjh.JSONGet("last_beat", ".")
	if err != nil {
		log.Printf("- Failed to get last beat: %v", err)
		return nil
	}

	var lastBeat HeartbeatBeat
	if err = json.Unmarshal(res.([]byte), &lastBeat); err != nil {
		panic(err)
	}
	return &lastBeat
}

// appendOrCreateArr will append an obj to key in path, or create objArr in key in path if it doesn't exist
func appendOrCreateArr(key string, path string, obj interface{}, objArr interface{}) error {
	// Check if key exists
	if _, err := rjh.JSONGet(key, path); err != nil {
		// Key doesn't exist, insert a new array with objArr
		if _, err := rjh.JSONSet(key, path, objArr); err != nil {
			return err
		}
	} else {
		// Key does exist, append obj to existing key
		_, err := rjh.JSONArrAppend(key, path, obj)
		return err
	}

	return nil
}
