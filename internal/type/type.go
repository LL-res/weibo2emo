package types

import "time"

type Post struct {
	Time     time.Time
	Content  string
	RawData  []string
	TimeData []int
	Likes    int
	Comments int
	Shares   int
	userID   string
	userType string
}
