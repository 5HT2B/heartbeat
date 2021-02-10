package main

import (
	"fmt"
	"strconv"
	"time"
)

func timeDifference(lastBeat string, t time.Time) string {
	lastBeatInt, err := strconv.ParseInt(lastBeat, 10, 64)

	if err != nil {
		return fmt.Sprintf("Could not convert '%s' to int64", lastBeat)
	}

	st := t.Sub(intToTime(lastBeatInt))

	return formattedTime(int(st.Seconds()))
}

func intToTime(int int64) time.Time {
	return time.Unix(int, 0)
}

func formattedTime(secondsIn int) string {
	hours := secondsIn / 3600
	minutes := (secondsIn / 60) - (60 * hours)
	seconds := secondsIn % 60
	str := fmt.Sprintf("%d hours, %d minutes, %d seconds", hours, minutes, seconds)
	return str
}
