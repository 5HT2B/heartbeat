package main

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

func ReadFileUnsafe(file string) string {
	content, err := ReadFile(file)

	if err != nil {
		log.Printf("- Failed to read '%s'", file)
	}

	return content
}

func ReadFile(file string) (string, error) {
	dat, err := ioutil.ReadFile(file)
	return string(dat), err
}

func WriteToFile(file string, content string) {
	data := []byte(content)
	err := ioutil.WriteFile(file, data, 0644)

	if err != nil {
		log.Printf("- Failed to read '%s'", file)
	}
}

// ReadLastBeatSafe returns the last beat and the longest period without a beat
func ReadLastBeatSafe() (int64, int64) {
	lastBeatStr, err := ReadFile("config/last_beat")

	if err != nil {
		return FixLastBeatFile()
	}

	split := strings.Split(lastBeatStr, ":")

	if len(split) != 2 {
		log.Printf("- config/last_beat file was %v, resetting", len(split))
		return FixLastBeatFile()
	}

	lastBeatInt, err := strconv.ParseInt(split[0], 10, 64)

	if err != nil {
		return FixLastBeatFile()
	}

	missingBeat, err := strconv.ParseInt(split[1], 10, 64)

	if err != nil {
		return FixLastBeatFile()
	}

	return lastBeatInt, missingBeat
}

func ReadGetRequestsSafe() int64 {
	totalVisitsStr, err := ReadFile("config/get_requests")

	if err != nil {
		return WriteGetRequestsFile(0)
	}

	split := strings.Split(totalVisitsStr, "\n")

	if len(split) != 2 {
		return WriteGetRequestsFile(0)
	}

	totalVisitsNew, err1 := strconv.ParseInt(split[1], 10, 64)

	if err1 != nil {
		return WriteGetRequestsFile(0)
	}

	return totalVisitsNew
}

func FixLastBeatFile() (int64, int64) {
	epoch := time.Now().Unix()
	WriteToFile("config/last_beat", strconv.FormatInt(epoch, 10)+":0")
	return epoch, 0
}

func WriteGetRequestsFile(int int64) int64 {
	WriteToFile("config/get_requests", "This page is only updated during a successful beat, so it may not always be up to date"+"\n"+strconv.FormatInt(int, 10))
	return int
}
