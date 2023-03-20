package types

import "time"

type Post struct {
	Time        time.Time
	Content     string
	RawData     []string
	TimeData    []int
	Likes       int
	Comments    int
	Shares      int
	UserID      string
	UserType    []string
	EmoKeyCount map[string]int //key en key, val count
}
type TypeName struct {
	Names map[string][]string `yaml:"names"`
}
