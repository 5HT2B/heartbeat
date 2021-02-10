package main

import (
	"io/ioutil"
	"log"
	"strconv"
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

func readLastBeatSafe() int64 {
	lastBeatStr, err := readFile("last_beat")

	if err != nil {
		writeToFile("last_beat", "0")
		return 0
	}

	lastBeatInt, err := strconv.ParseInt(lastBeatStr, 10, 64)

	if err != nil {
		writeToFile("last_beat", "0")
		return 0
	}

	return lastBeatInt
}
