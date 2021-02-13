package main

import (
	"fmt"
	"time"
)

func TimeDifference(lastBeat int64, t time.Time) string {
	st := t.Sub(IntToTime(lastBeat))
	return FormattedTime(int(st.Seconds()))
}

func IntToTime(int int64) time.Time {
	return time.Unix(int, 0)
}

func FormattedTime(secondsIn int) string {
	hours := secondsIn / 3600
	minutes := (secondsIn / 60) - (60 * hours)
	seconds := secondsIn % 60
	str := fmt.Sprintf("%d hours, %d minutes, %d seconds", hours, minutes, seconds)
	return str
}
