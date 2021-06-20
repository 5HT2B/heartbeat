package main

import (
	"fmt"
	"time"
)

func TimeDifference(lastBeat int64, t time.Time) string {
	st := t.Sub(time.Unix(lastBeat, 0))
	return FormattedTime(int(st.Seconds()))
}

func FormattedTime(secondsIn int) string {
	hours := secondsIn / 3600
	minutes := (secondsIn / 60) - (60 * hours)
	seconds := secondsIn % 60
	str := fmt.Sprintf("%d hours, %d minutes, %d seconds", hours, minutes, seconds)
	return str
}
