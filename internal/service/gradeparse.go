package service

import (
	"encoding/csv"
	"github.com/emomo/weibo2emo/internal/tools"
	"github.com/schollz/progressbar/v3"
	"io"
	"log"
	"os"
)

type GradeParser struct {
	Grade int
	Data  [][]string
	Heads []string
	Times []string
}

func NewGradeParser(grade int) *GradeParser {
	return &GradeParser{
		Grade: grade,
		Data:  make([][]string, 0),
		Heads: make([]string, 0),
		Times: make([]string, 0),
	}
}
func (g *GradeParser) LoadGradeData(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	//reader := transform.NewReader(f, simplifiedchinese.GBK.NewDecoder())
	r := csv.NewReader(f)
	r.LazyQuotes = true
	isHead := true
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("E! 读取分级数据源时有一行数据读取失败，但我们接着读", err)
			continue
		}
		if isHead {
			g.Heads = row
			isHead = false
			continue
		}
		g.Data = append(g.Data, row[1:])
		g.Times = append(g.Times, row[0])

	}

	return nil
}
func (g *GradeParser) ExportResult(path string, trans func(float64) float64) error {
	bar := progressbar.Default(int64(len(g.Data)), "开始计算并导出结果")
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

	maxSlice := make([]string, len(g.Heads)-1)
	for i, v := range g.Data {
		if i == 0 {
			copy(maxSlice, v)
			continue
		}
		for index, score := range v {
			maxSlice[index] = tools.GetMaxS(score, maxSlice[index])
		}
	}
	result := make([][]string, 0, len(g.Data))
	result = append(result, g.Heads)
	for i, v := range g.Data {
		grades, err := tools.GetGrades(v, maxSlice, g.Grade, trans)
		if err != nil {
			log.Println("获取等级失败，时间 :", g.Times[i], err)
			result = append(result, []string{g.Times[i]})
			continue
		}
		temp := []string{g.Times[i]}
		temp = append(temp, grades...)
		result = append(result, temp)
		bar.Add(1)
	}
	w.WriteAll(result)
	w.Flush()
	return nil
}
