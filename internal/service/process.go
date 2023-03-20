package service

import (
	"encoding/csv"
	"fmt"
	"github.com/emomo/weibo2emo/database"
	"github.com/emomo/weibo2emo/internal/tools"
	types "github.com/emomo/weibo2emo/internal/type"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/yanyiwu/gojieba"
)

type TypeMap map[string][]string

type Processor struct {
	Dictionary  map[string][]string // key : emo kind , val :words
	Posts       []types.Post
	Slicer      *gojieba.Jieba
	CountMap    map[string]int
	EmoKey      []string
	TypeMap     TypeMap //key :关键的用户id , val : 这个id所属的类型，可能是多种
	TypeNum     int
	ToDB        bool
	Concurrency int
}

func (p *Processor) FillTypeMap(path string) error {
	f, err := os.Open("./conf/catagory.yaml")
	if err != nil {
		log.Println("无法打开类型名称配置文件", err)
		return err
	}
	defer f.Close()
	bytes, err := io.ReadAll(f)
	if err != nil {
		log.Println("无法读取类型名称配置文件", err)
		return err
	}
	names := types.TypeName{}           // key 种类 val 包含的词
	namesMap := make(map[string]string) //key 包含的词，val 这个词的类型
	err = yaml.Unmarshal(bytes, &names)
	if err != nil {
		log.Println("unmarshal err", err)
		return err
	}
	p.TypeNum = len(names.Names) + 1
	for k, ns := range names.Names {
		fillMap(namesMap, k, ns)
	}

	file, err := os.Open(path)
	if err != nil {
		log.Println("无法打开 名称 id 表", err)
		return err
	}
	defer file.Close()
	r, err := tools.DetectCSV(file)
	if err != nil {
		log.Println("WARN ! 编码格式检测失败，使用utf-8开始解码")
	}
	head := true
	for num := 1; ; num++ {
		row, err := r.Read()
		if head {
			head = false
			continue
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("读取 名称 id 表%d行失败%v", num, err)
			continue
		}
		for k, ns := range names.Names { // k : media,municiplity,authority ,ns : chinese names
			for _, n := range ns { // n : name
				if strings.Contains(row[0], n) { // row[0] : content
					p.TypeMap.insert(row[1], k) // row[1] : id
				}
			}
		}
	}
	return nil
}
func NewProcessor() *Processor {
	dictDir := path.Join(filepath.Dir(os.Args[0]), "dict")
	jiebaPath := path.Join(dictDir, "jieba.dict.utf8")
	hmmPath := path.Join(dictDir, "hmm_model.utf8")
	userPath := path.Join(dictDir, "user.dict.utf8")
	idfPath := path.Join(dictDir, "idf.utf8")
	stopPath := path.Join(dictDir, "stop_words.utf8")
	return &Processor{
		Dictionary: make(map[string][]string),
		Posts:      make([]types.Post, 0),
		Slicer:     gojieba.NewJieba(jiebaPath, hmmPath, userPath, idfPath, stopPath),
		CountMap:   make(map[string]int),
		EmoKey:     make([]string, 0),
		TypeMap:    NewTypeMap(),
	}
}

func (p *Processor) LoadDictionary(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	r, err := tools.DetectCSV(f)
	if err != nil {
		log.Println("WARN ! 编码格式检测失败，使用utf-8开始解码")
	}
	r.LazyQuotes = true
	for {
		rawRow, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("E! 读取词库时有一行数据读取失败", err)
			continue
		}
		row := make([]string, 0)
		for _, v := range rawRow {
			if v != "" {
				row = append(row, v)
			}
		}
		p.EmoKey = append(p.EmoKey, row[0])
		p.Dictionary[row[0]] = row[1:]
		for i := 1; i < len(row); i++ {
			p.Slicer.AddWord(row[i])
		}
	}

	return nil
}
func (p *Processor) LoadPosts(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	r, err := tools.DetectCSV(f)
	if err != nil {
		log.Println("WARN ! 编码格式检测失败，使用utf-8开始解码")
	}
	r.LazyQuotes = true
	isheader := true
	for num := 1; ; num++ {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if isheader {
			isheader = false
			continue
		}
		if err != nil {
			log.Printf("E! 读取博文时%d行读取失败 %v", num, err)
			continue
		}
		t, err := time.Parse("01月02日 15:04", row[2])
		t = t.AddDate(2020, 0, 0)
		if err != nil {
			log.Printf("E! 读取博文时%d行数据时间转换失败 %v", num, err)
			//log.Println("E! 读取博文时有一行数据时间转换失败，但我们接着读", err)
			continue
		}
		likes, _ := strconv.Atoi(row[3])
		comments, _ := strconv.Atoi(row[4])
		shares, _ := strconv.Atoi(row[5])
		post := types.Post{
			Time:        t,
			Content:     row[1],
			RawData:     row,
			Likes:       likes,
			Comments:    comments,
			Shares:      shares,
			UserID:      row[0],
			UserType:    p.TypeMap[row[0]],
			EmoKeyCount: make(map[string]int, 0),
		}
		p.Posts = append(p.Posts, post)
	}
	return nil
}
func (p *Processor) ExportResult(path string) error {
	bar := progressbar.Default(int64(len(p.Posts)), "开始计算并导出结果")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString("\xEF\xBB\xBF")
	if err != nil {
		log.Println("E ! 写入Bomm头错误", err)
	}
	w := csv.NewWriter(f)

	data := []string{"博主id", "博文", "发布时间", "点赞数", "评论数", "转发数", "分词结果"}
	data = append(data, p.EmoKey...)
	err = w.Write(data)
	if err != nil {
		log.Println("E ! 写入结果数据集时首行错误", err)
	}
	w.Flush()
	dataToFlush := make([][]string, len(p.Posts))
	group := sync.WaitGroup{}
	barLock := sync.Mutex{}
	limitChan := make(chan struct{}, p.Concurrency)
	group.Add(len(p.Posts))
	for i, post := range p.Posts {
		limitChan <- struct{}{}
		go func(i int, post types.Post) {
			defer group.Done()
			countMap := make(map[string]int)
			slices := p.Slicer.Cut(tools.RemoveName(post.Content), true)
			for _, s := range slices {
				existKeys := tools.CheckExistanceInMap(s, p.Dictionary)
				for _, keyEmo := range existKeys {
					countMap[keyEmo]++
				}
			}
			toWriteDate := post.RawData
			toWriteDate = append(toWriteDate, tools.RemoveMark(slices))
			for _, v := range p.EmoKey {
				toWriteDate = append(toWriteDate, strconv.Itoa(countMap[v]))
				p.Posts[i].EmoKeyCount[tools.ENEmoKey(v)] = countMap[v]
				p.Posts[i].TimeData = append(p.Posts[i].TimeData, countMap[v])
			}
			p.Posts[i].TimeData = append(p.Posts[i].TimeData, post.Likes, post.Comments, post.Shares)
			dataToFlush[i] = toWriteDate
			//err = w.Write(toWriteDate)
			barLock.Lock()
			bar.Add(1)
			barLock.Unlock()
			<-limitChan
		}(i, post)
	}
	group.Wait()
	w.WriteAll(dataToFlush)
	w.Flush()
	return p.FlushToDB()
}
func (p *Processor) FlushToDB() error {
	if !p.ToDB {
		return nil
	}
	db, err := database.NewPostDB()
	if err != nil {
		return err
	}
	err = db.NewPostsTable(p.EmoKey)
	if err != nil {
		return err
	}
	db.InsertData(p.Posts, p.EmoKey)
	return nil

}
func (p *Processor) ExportResultByTime(path string) error {
	bar := progressbar.Default(int64(len(p.Posts)), "开始计算并导出按时间聚合的结果")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString("\xEF\xBB\xBF")
	if err != nil {
		log.Println("E ! 写入Bomm头错误", err)
	}
	w := csv.NewWriter(f)
	tmap := make(map[string][]int)       //key : 时间与类型拼接 val 各项指标的值，如情绪，点赞转发数等
	timeSet := make(map[string]struct{}) // 时间
	for i := 0; i < len(p.Posts); i++ {
		timeKey := p.Posts[i].Time.Format("01月02日 15 时")
		timeSet[timeKey] = struct{}{}
		keys := make([]string, 0) //一条博文的userType可能有多种

		if p.Posts[i].UserType == nil {
			keys = append(keys, fmt.Sprintf("%s$A-normal", timeKey))
		} else {
			for _, v := range p.Posts[i].UserType {
				keys = append(keys, fmt.Sprintf("%s$%s", timeKey, v))
			}
		}
		for _, key := range keys {
			if _, ok := tmap[key]; ok {
				tools.AddTwoSlices(tmap[key], p.Posts[i].TimeData)
			} else {
				tt := make([]int, len(p.Posts[i].TimeData))
				copy(tt, p.Posts[i].TimeData)
				tmap[key] = tt
			}
		}
		bar.Add(1)
	}
	dataToWrite := make([][]string, 0, len(timeSet))
	header := []string{"日期"}
	three := []string{"点赞数", "评论数", "转发数"}
	for i := 1; i <= p.TypeNum; i++ {
		for _, v := range append(p.EmoKey, three...) { //把情绪跟那几个数后面都加一个第几类的后缀
			header = append(header, fmt.Sprintf("%s-第%d类", v, i))
		}
	}
	dataToWrite = append(dataToWrite, header)
	sortTemp := make([]string, 0, len(tmap))
	for k := range tmap {
		sortTemp = append(sortTemp, k)
	}
	sortTime := make([]string, 0, len(timeSet))
	for k := range timeSet {
		sortTime = append(sortTime, k)
	}
	sort.Strings(sortTime) // 01 01
	sort.Strings(sortTemp) // 01 01 $ A

	var ptr int
	for _, v := range sortTime {
		singleData := make([]string, (p.TypeNum)*(len(p.EmoKey)+3)+1) //有几类就会有几个维度，并且还会有一个最开始的时间
		singleData[0] = v
		for ; ptr < len(sortTemp); ptr++ {
			if !strings.HasPrefix(sortTemp[ptr], v) {
				break //如果完全键的时间不在是传入的时间，则进行下一个时间的统计
			}
			level := extractLevel(sortTemp[ptr])
			singleWrite := tools.ConvIs2Ss(tmap[sortTemp[ptr]])
			for i, num := range singleWrite {
				singleData[1+(level-1)*(len(p.EmoKey)+3)+i] = num
			}
		}
		dataToWrite = append(dataToWrite, singleData)
	}
	err = w.WriteAll(dataToWrite)
	if err != nil {
		log.Println("写入时间结果数据集时错误 ", err)
	}
	w.Flush()
	return nil

}
func (p *Processor) resetCountMap() {
	p.CountMap = make(map[string]int)
}
func fillMap(namesMap map[string]string, val string, keys []string) {
	for _, k := range keys {
		namesMap[k] = val
	}
}
func (t TypeMap) insert(key string, val string) {
	if t[key] == nil {
		t[key] = make([]string, 0)
	}
	for _, v := range t[key] {
		if v == val {
			return
		}
	}
	t[key] = append(t[key], val)
}
func NewTypeMap() TypeMap {
	t := make(map[string][]string)
	return t
}

// given 01月01日 15时$A-normal output 1
func extractLevel(s string) int {
	strs := strings.Split(s, "$")
	return int(strs[1][0] - 'A' + 1)
}
