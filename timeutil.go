package main

import (
	"fmt"
	"time"
)

func timeDifference(lastBeat int64, t time.Time) string {
	st := t.Sub(intToTime(lastBeat))

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
