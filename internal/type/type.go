package types

import "time"

type Post struct {
	Time     time.Time
	Content  string
	RawData  []string
	TimeData []int
}
