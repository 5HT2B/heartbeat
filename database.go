package main

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
	"log"
	"time"
)

var (
	ctx = context.Background()
	rdb *redis.Client   // set in main()
	rjh *rejson.Handler // set in main()

	redisAddr = "localhost:6379" // set by .env
	redisPass = ""               // set by .env

	lastBeat         int64              = -1
	missingBeat      int64              = -1
	totalBeats       int64              = -1
	totalVisits      int64              = -1
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

func CurrentMissingBeat() int64 {
	// todo: this needs to work for multiple devices
	lastBeatTmp := GetLastBeat() // todo: rename
	diff := time.Now().Unix() - lastBeatTmp.Timestamp

	if diff > heartbeatStats.LongestMissingBeat {
		return diff
	}
	return heartbeatStats.LongestMissingBeat
}

func UpdateTotalVisits() {
	// TODO: unfinished
	if totalVisits == -1 || heartbeatStats == nil {
		if res, err := rjh.JSONGet("stats", "."); err != nil {
			if _, err := rjh.JSONSet("stats", ".", defaultHeartbeatStats); err != nil {
				panic(err)
			}
		} else {
			var stats HeartbeatStats
			if err = json.Unmarshal(res.([]byte), &stats); err != nil {
				panic(err)
			}

			stats.TotalVisits += 1

		}
	}
}

func UpdateDevices(lastBeat *HeartbeatBeat, increment bool) {
	//for _, device := range *heartbeatDevices {
	//
	//}
	// TODO do this
}

// SetLastBeat will save the last beat and insert a new HeartbeatBeat into beats
func SetLastBeat(deviceName string, timestamp int64) error {
	heartbeatStats.TotalBeats += 1

	lastBeat := HeartbeatBeat{deviceName, timestamp}
	if _, err := rjh.JSONSet("last_beat", ".", lastBeat); err != nil {
		return err
	}

	lastBeatArr := []HeartbeatBeat{lastBeat}
	err := appendOrCreateArr("beats", ".", lastBeat, lastBeatArr)

	return err
}

// GetLastBeat will get the last beat, and return nilLastBeat if there was an error retrieving it (likely no beat)
func GetLastBeat() *HeartbeatBeat {
	if rjh == nil {
		time.Sleep(1 * time.Second)
		return GetLastBeat()
	}

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
