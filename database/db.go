package database

import (
	"fmt"
	"github.com/emomo/weibo2emo/internal/tools"
	types "github.com/emomo/weibo2emo/internal/type"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"reflect"
	"strings"
	"time"
)

type PostDB struct {
	db *gorm.DB
}

func NewPostDB() (*PostDB, error) {
	dsn := "root:root@tcp(127.0.0.1:3306)/idb?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetConnMaxLifetime(time.Hour)
	db.PrepareStmt = false
	db.Logger = logger.Default.LogMode(logger.Error)
	return &PostDB{db: db}, nil
}

func (p *PostDB) NewPostsTable(emoKeys []string) error {
	if p.db.Migrator().HasTable("posts") {
		err := p.db.Migrator().DropTable("posts")
		if err != nil {
			log.Println(err)
			return err
		}
	}
	fields := []reflect.StructField{
		{
			Name: "UserID",
			Type: reflect.TypeOf(uint(1)),
			Tag:  `gorm:"association_autoupdate:false"`,
		},
		{
			Name: "UserType",
			Type: reflect.TypeOf(uint(1)),
			Tag:  `gorm:"association_autoupdate:false"`,
		},
		{
			Name: "Time",
			Type: reflect.TypeOf(time.Time{}),
			Tag:  `gorm:"association_autoupdate:false"`,
		},
		//{
		//	Name: "Content",
		//	Type: reflect.TypeOf(""),
		//},
		{
			Name: "Likes",
			Type: reflect.TypeOf(uint(1)),
			Tag:  `gorm:"association_autoupdate:false"`,
		},
		{
			Name: "Comments",
			Type: reflect.TypeOf(uint(1)),
			Tag:  `gorm:"association_autoupdate:false"`,
		},
		{
			Name: "Shares",
			Type: reflect.TypeOf(uint(1)),
			Tag:  `gorm:"association_autoupdate:false"`,
		},
	}
	for _, key := range emoKeys {
		fields = append(fields, reflect.StructField{
			Name: tools.ENEmoKey(key),
			Type: reflect.TypeOf(uint(1)),
			Tag:  `gorm:"association_autoupdate:false"`,
		})
	}
	//log.Println(fields)
	t := reflect.StructOf(fields)
	table := reflect.New(t).Interface()
	err := p.db.Table("posts").AutoMigrate(table)
	if err != nil {
		log.Println(err)
	}
	return err
}
func (p *PostDB) InsertData(posts []types.Post, emoKeys []string) {
	vals := make([]map[string]interface{}, 0)
	for _, post := range posts {
		singleRecord := map[string]interface{}{
			"user_id":   post.UserID,
			"user_type": BinaryType(post),
			"time":      post.Time,
			//"content":   post.Content,
			"likes":    post.Likes,
			"comments": post.Comments,
			"shares":   post.Shares,
		}
		for _, chkey := range emoKeys {
			enkey := tools.ENEmoKey(chkey)
			singleRecord[enkey] = post.EmoKeyCount[enkey]
		}
		vals = append(vals, singleRecord)
	}
	if len(vals) == 0 {
		return
	}
	fmt.Println("batch size : ", 65535/len(vals[0]))
	err := p.db.Table("posts").CreateInBatches(vals, 65535/len(vals[0])).Error //.Create(vals)
	if err != nil {
		log.Println(err)
	}
}
func BinaryType(strType types.Post) int {
	if strType.UserType == nil {
		return 1
	}
	res := 0
	for _, t := range strType.UserType {
		strs := strings.Split(t, "-")
		res |= 1 << (strs[0][0] - 'A')
	}
	return res
}
