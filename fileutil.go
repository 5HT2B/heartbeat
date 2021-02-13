package main

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
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

// Returns the last beat and the longest period without a beat
func ReadLastBeatSafe() (int64, int64) {
	lastBeatStr, err := ReadFile("last_beat")

	if err != nil {
		return FixLastBeatFile()
	}

	split := strings.Split(lastBeatStr, ":")

	if len(split) != 2 {
		log.Printf("- last_beat file was %v, resetting", len(split))
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

func FixLastBeatFile() (int64, int64) {
	WriteToFile("last_beat", "0:0")
	return 0, 0
}
