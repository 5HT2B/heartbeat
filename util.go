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

func TimeDifference(lastBeat int64, t time.Time) string {
	st := t.Sub(time.Unix(lastBeat, 0))
	return FormattedTime(int64(st.Seconds()))
}

func FormattedTime(secondsIn int64) string {
	hours := secondsIn / 3600
	minutes := (secondsIn / 60) - (60 * hours)
	seconds := secondsIn % 60
	str := fmt.Sprintf("%d hours, %d minutes, %d seconds", hours, minutes, seconds)
	return str
}

func FormattedNum(num int64) string {
	return printer.Sprintf("%d", num)
}
