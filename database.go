package main

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
	"log"
	"time"
)

var (
	rdb *redis.Client   // set in main()
	rjh *rejson.Handler // set in main()

	redisAddr = "localhost:6379" // set by .env
	redisPass = ""               // set by .env

	heartbeatStats   *HeartbeatStats    = nil
	heartbeatDevices *[]HeartbeatDevice = nil
)

func SetupDatabase() (*redis.Client, *rejson.Handler) {
	rh := rejson.NewReJSONHandler()
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPass,
	})

	rh.SetGoRedisClient(client)
	return client, rh
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

func UpdateTotalVisits() {
	// TODO: unfinished
	//if totalVisits == -1 || heartbeatStats == nil {
	//	if res, err := rjh.JSONGet("stats", "."); err != nil {
	//		if _, err := rjh.JSONSet("stats", ".", defaultHeartbeatStats); err != nil {
	//			panic(err)
	//		}
	//	} else {
	//		var stats HeartbeatStats
	//		if err = json.Unmarshal(res.([]byte), &stats); err != nil {
	//			panic(err)
	//		}
	//
	//		stats.TotalVisits += 1
	//
	//	}
	//}
}

// LastSeen will return the formatted date of the last timestamp received from a beat
func LastSeen() string {
	lastBeat := GetLastBeat()
	return time.Unix(lastBeat.Timestamp, 0).Format(timeFormat)
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

// UpdateUptime will update the uptime statistics
func UpdateUptime() {
	now := time.Now().Unix()
	diff := now - uptimeTimer

	uptimeTimer = now - 1
	heartbeatStats.TotalUptime += diff
}

// UpdateDevice will update the LastBeat of a device
func UpdateDevice(beat HeartbeatBeat) {
	for n, device := range *heartbeatDevices {
		if device.DeviceName == beat.DeviceName {
			diff := time.Now().Unix() - device.LastBeat.Timestamp
			if diff > device.LongestMissingBeat {
				device.LongestMissingBeat = diff
			}
			if device.LongestMissingBeat > heartbeatStats.LongestMissingBeat {
				heartbeatStats.LongestMissingBeat = device.LongestMissingBeat
			}

			device.LastBeat = beat
			device.TotalBeats += 1

			(*heartbeatDevices)[n] = device
			heartbeatStats.TotalBeats += 1
			break // name should only ever be matching once
		}
	}
}

// UpdateLastBeat will save the last beat and insert a new HeartbeatBeat into beats
func UpdateLastBeat(deviceName string, timestamp int64) error {
	lastBeat := HeartbeatBeat{deviceName, timestamp}
	UpdateDevice(lastBeat)

	if _, err := rjh.JSONSet("last_beat", ".", lastBeat); err != nil {
		return err
	}

	lastBeatArr := []HeartbeatBeat{lastBeat}
	err := appendOrCreateArr("beats", ".", lastBeat, lastBeatArr)

	return err
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
