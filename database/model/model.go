package model

import "time"

type Post struct {
	Time     time.Time
	Likes    int
	Comments int
	Shares   int
	UserID   int //主键
	UserType int
	items    map[string]int //二进制表示，因为可能有多种类型，1代表普通 2代表B类，4代表C类，8代表D类
}
