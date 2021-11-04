package main

import "fmt"

//
// Redis JSONPath structure is documented in DOCS.md
//

var (
	defaultHeartbeatStats   = &HeartbeatStats{"Never", "0", 0, 0, 0, 0}
	defaultHeartbeatDevices = &[]HeartbeatDevice{}
)

// HeartbeatInfo is the human-readable info presented on the main page
type HeartbeatInfo struct {
	LastSeen       string `json:"last_seen"`
	TimeDifference string `json:"time_difference"`
	MissingBeat    string `json:"missing_beat"`
	TotalBeats     string `json:"total_beats"`
}

// HeartbeatBeat is the current last beat
type HeartbeatBeat struct {
	DeviceName string `json:"device_name"`
	Timestamp  int64  `json:"timestamp"`
}

// HeartbeatDevice is used in an array of recognized devices
type HeartbeatDevice struct {
	DeviceName         string        `json:"device_name"`
	LastBeat           HeartbeatBeat `json:"last_beat"`
	TotalBeats         int64         `json:"total_beats"`
	LongestMissingBeat int64         `json:"longest_missing_beat"`
}

// HeartbeatStats are the global stats for a heartbeat server
type HeartbeatStats struct {
	LastBeatFormatted   string `json:"last_beat_formatted,omitempty"`   // handled by UpdateLastBeatFmt, UpdateLastBeatFmtV and LongestAbsence
	TotalBeatsFormatted string `json:"total_beats_formatted,omitempty"` // handled by UpdateTotalBeats
	TotalVisits         int64  `json:"total_visits"`                    // handled by HandleSuccessfulBeat
	TotalUptime         int64  `json:"total_uptime"`                    // handled by UpdateUptime
	TotalBeats          int64  `json:"total_beats"`                     // handled by UpdateDevice
	LongestMissingBeat  int64  `json:"longest_missing_beat"`            // handled by LongestAbsence
}

func (lb HeartbeatBeat) String() string {
	return fmt.Sprintf("HeartbeatBeat<%s, %v>", lb.DeviceName, lb.Timestamp)
}

func (s HeartbeatDevice) String() string {
	return fmt.Sprintf("HeartbeatDevice<%s, %v, %v>", s.DeviceName, s.TotalBeats, s.LongestMissingBeat)
}

func (s HeartbeatStats) String() string {
	return fmt.Sprintf("HeartbeatStats<%v, %v, %v, %v>", s.TotalVisits, s.TotalUptime, s.TotalBeats, s.LongestMissingBeat)
}
