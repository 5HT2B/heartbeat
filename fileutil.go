package main

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func readFileUnsafe(file string) string {
	content, err := readFile(file)

	if err != nil {
		log.Printf("- Failed to read '%s'", file)
	}

	return content
}

func readFile(file string) (string, error) {
	dat, err := ioutil.ReadFile(file)
	return string(dat), err
}

func writeToFile(file string, content string) {
	data := []byte(content)
	err := ioutil.WriteFile(file, data, 0644)

	if err != nil {
		log.Printf("- Failed to read '%s'", file)
	}
}

// Returns the last beat and the longest period without a beat
func readLastBeatSafe() (int64, int64) {
	lastBeatStr, err := readFile("last_beat")

	if err != nil {
		return fixLastBeatFile()
	}

	split := strings.Split(lastBeatStr, ":")

	if len(split) != 2 {
		log.Printf("- last_beat file was %v, resetting", len(split))
		return fixLastBeatFile()
	}

	lastBeatInt, err := strconv.ParseInt(split[0], 10, 64)

	if err != nil {
		return fixLastBeatFile()
	}

	missingBeat, err := strconv.ParseInt(split[1], 10, 64)

	if err != nil {
		return fixLastBeatFile()
	}

	return lastBeatInt, missingBeat
}

func fixLastBeatFile() (int64, int64) {
	writeToFile("last_beat", "0:0")
	return 0, 0
}
