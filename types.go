package main

import "fmt"

//
// Redis JSONPath structure is documented in DOCS.md
//

var (
	defaultHeartbeatStats = HeartbeatStats{0, 0, 0, 0}
)

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
	TotalVisits        int64 `json:"total_visits"`
	TotalUptime        int64 `json:"total_uptime"`
	TotalBeats         int64 `json:"total_beats"`          // handled by updateDevice
	LongestMissingBeat int64 `json:"longest_missing_beat"` // handled by updateDevice
}

func (lb HeartbeatBeat) String() string {
	return fmt.Sprintf("HeartbeatBeat<%s, %v>", lb.DeviceName, lb.Timestamp)
}

func (s HeartbeatDevice) String() string {
	return fmt.Sprintf("HeartbeatDevice<%s, %v, %v>", s.DeviceName, s.TotalBeats, s.LongestMissingBeat)
}

func (s HeartbeatStats) String() string {
	return fmt.Sprintf("HeartbeatStats<%v, %v, %v>", s.TotalVisits, s.TotalBeats, s.LongestMissingBeat)
}
