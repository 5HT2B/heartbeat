package main

import (
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"time"
)

var (
	printer = message.NewPrinter(language.English)
)

// TimeDifference will return the FormattedTime of time.Unix().Now() - lastBeat
func TimeDifference(lastBeat int64, t time.Time) string {
	st := t.Sub(time.Unix(lastBeat, 0))
	return FormattedTime(int64(st.Seconds()))
}

// FormattedTime will turn seconds into a pretty time representation
func FormattedTime(secondsIn int64) string {
	hours := secondsIn / 3600
	minutes := (secondsIn / 60) - (60 * hours)
	seconds := secondsIn % 60
	str := fmt.Sprintf("%d hours, %d minutes, %d seconds", hours, minutes, seconds)
	return str
}

// FormattedNum will insert commas as necessary in large numbers
func FormattedNum(num int64) string {
	return printer.Sprintf("%d", num)
}
