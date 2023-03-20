package database

import (
	types "github.com/emomo/weibo2emo/internal/type"
	"testing"
	"time"
)

func TestConnect(t *testing.T) {
	db, err := NewPostDB()
	if err != nil {
		panic(err)
	}
	err = db.NewPostsTable([]string{"失望(NJ)", "慌(NI)", "死亡词", "提及政府当局", "nb(NB)", "hh(HH)", "gg(GG)", "jj(JJ)", "kk(KK)", "ll(LL)", "zz(ZZ)", "xx(XX)"})
	if err != nil {
		panic(err)
	}
	ts, err := time.Parse("01月02日 15:04", "01月01日 21:26")
	ts = ts.AddDate(2020, 0, 0)
	post := types.Post{
		Time:        ts,
		Content:     "我是谁",
		Likes:       1,
		Comments:    4,
		Shares:      5,
		UserID:      "1454",
		UserType:    []string{"C-municipality", "B-media"},
		EmoKeyCount: map[string]int{"Death": 4, "NJ": 6, "NB": 7, "XX": 9, "JJ": 10},
	}
	posts := []types.Post{}
	for i := 0; i < 1200000; i++ {
		posts = append(posts, post)
	}
	//posts := []types.Post{
	//	{
	//		Time:        ts,
	//		Content:     "我是谁",
	//		Likes:       1,
	//		Comments:    2,
	//		Shares:      3,
	//		UserID:      "1232",
	//		UserType:    nil,
	//		EmoKeyCount: nil,
	//	},
	//	{
	//		Time:        ts,
	//		Content:     "我是谁",
	//		Likes:       1,
	//		Comments:    4,
	//		Shares:      5,
	//		UserID:      "1454",
	//		UserType:    []string{"C-municipality", "B-media"},
	//		EmoKeyCount: map[string]int{"Death": 4, "NJ": 6},
	//	},
	//}
	db.InsertData(posts, []string{"失望(NJ)", "慌(NI)", "死亡词", "提及政府当局", "nb(NB)", "hh(HH)", "gg(GG)", "jj(JJ)", "kk(KK)", "ll(LL)", "zz(ZZ)", "xx(XX)"})

}
