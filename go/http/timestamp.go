package http

import "time"

const (
	httpTimestampHeader      = "X-Termly-Timestamp"
	httpTimestampValueFormat = "20060102T150405Z"
)

func newTimestamp() string {
	return time.Now().UTC().Format(httpTimestampValueFormat)
}
