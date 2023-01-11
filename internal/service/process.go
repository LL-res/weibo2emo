package service

import (
	"encoding/csv"
	"github.com/emomo/weibo2emo/internal/tools"
	types "github.com/emomo/weibo2emo/internal/type"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/yanyiwu/gojieba"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type Processor struct {
	Dictionary map[string][]string // key : emo kind , val :words
	Posts      []types.Post
	Slicer     *gojieba.Jieba
	CountMap   map[string]int
	EmoKey     []string
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
	}
}

func (p *Processor) LoadDictionary(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	reader := transform.NewReader(f, simplifiedchinese.GBK.NewDecoder())
	r := csv.NewReader(reader)
	r.LazyQuotes = true
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("E! 读取词库时有一行数据读取失败，但我们接着读", err)
			continue
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
	reader := transform.NewReader(f, simplifiedchinese.GBK.NewDecoder())
	r := csv.NewReader(reader)
	r.LazyQuotes = true
	isheader := true
	for num := 1; ; num++ {
		if isheader {
			isheader = false
			continue
		}
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("E! 读取博文时%d行读取失败 %v", num, err)
			continue
		}
		t, err := time.Parse("01月02日 15:04", row[2])

		if err != nil {
			log.Printf("E! 读取博文时%d行数据时间转换失败 %v", num, err)
			//log.Println("E! 读取博文时有一行数据时间转换失败，但我们接着读", err)
			continue
		}
		likes, _ := strconv.Atoi(row[3])
		comments, _ := strconv.Atoi(row[4])
		shares, _ := strconv.Atoi(row[5])
		post := types.Post{
			Time:     t,
			Content:  row[1],
			RawData:  row,
			Likes:    likes,
			Comments: comments,
			Shares:   shares,
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
	limitChan := make(chan struct{}, 5)
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
	tmap := make(map[string][]int) // key date ,val
	for i := 0; i < len(p.Posts); i++ {
		key := p.Posts[i].Time.Format("01月02日 15 时")
		if _, ok := tmap[key]; ok {
			tools.AddTwoSlices(tmap[key], p.Posts[i].TimeData)
		} else {
			tt := make([]int, len(p.Posts[i].TimeData))
			copy(tt, p.Posts[i].TimeData)
			tmap[key] = tt
		}
		bar.Add(1)
	}
	dataToWrite := make([][]string, 0, len(tmap))
	header := []string{"日期"}
	header = append(header, p.EmoKey...)
	header = append(header, "点赞数", "评论数", "转发数")
	dataToWrite = append(dataToWrite, header)
	sortTemp := make([]string, 0, len(tmap))
	for k := range tmap {
		sortTemp = append(sortTemp, k)
	}
	sort.Strings(sortTemp)
	for _, v := range sortTemp {
		singleData := []string{v}
		singleData = append(singleData, tools.ConvIs2Ss(tmap[v])...)
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
