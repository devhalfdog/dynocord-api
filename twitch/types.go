package twitch

import "time"

type Segment struct {
	URI           string
	TotalDuration time.Duration
}
