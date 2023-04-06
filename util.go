package main

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
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

	units := make([]string, 0)
	if hours != 0 {
		units = append(units, joinIntAndStr(hours, "hour"))
	}
	if minutes != 0 {
		units = append(units, joinIntAndStr(minutes, "minute"))
	}
	if seconds != 0 || (hours == 0 && minutes == 0) {
		units = append(units, joinIntAndStr(seconds, "second"))
	}

	return strings.Join(units, ", ")
}

// FormattedNum will insert commas as necessary in large numbers
func FormattedNum(num int64) string {
	return printer.Sprintf("%d", num)
}

func joinIntAndStr(int int64, str string) string {
	plural := "s"
	if int == 1 {
		plural = ""
	}
	return fmt.Sprintf("%s %s%s", FormattedNum(int), str, plural)
}
