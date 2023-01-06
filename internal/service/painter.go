package service

import (
	"encoding/csv"
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Painter struct {
	lines map[string][]float64
	base  []string
	emos  []string
	size  []int
}

func NewPainter(width, height int) *Painter {
	return &Painter{
		lines: make(map[string][]float64),
		base:  make([]string, 0),
		emos:  make([]string, 0),
		size:  []int{width, height},
	}
}

func (p *Painter) Load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
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
			for i := 1; i < len(row); i++ {
				p.lines[row[i]] = make([]float64, 0)
				p.emos = append(p.emos, row[i])

			}
			isHead = false
			continue
		}
		for i, v := range row {
			if i == 0 {
				p.base = append(p.base, v)
				continue
			}
			t, _ := strconv.ParseFloat(v, 64)
			p.lines[p.emos[i-1]] = append(p.lines[p.emos[i-1]], t)
		}

	}
	for i := range p.base {
		t, _ := time.Parse("01月02日 15 时", p.base[i])
		p.base[i] = t.Format("01M\n02D\n15h")
	}
	return nil
}
func (p *Painter) Paint() error {
	temp := make([]interface{}, 0)
	temp1n := make([]interface{}, 0)
	for k := range p.lines {
		if strings.HasPrefix(k, "ln") {
			temp1n = append(temp1n, k, p.drawPoints(k))
		}
		temp = append(temp, k, p.drawPoints(k))
	}
	err := p.paint(temp)
	if err != nil {
		return err
	}
	err = p.paint(temp1n)
	if err != nil {
		return err
	}
	return nil
}
func (p *Painter) drawPoints(key string) plotter.XYs {
	points := make(plotter.XYs, len(p.lines[key]))
	for i := range points {
		points[i].X = float64(i)
		points[i].Y = p.lines[key][i]
	}
	return points
}
func (p *Painter) paint(points []interface{}) error {
	r := plot.New()
	r.X.Label.Text = "度量值"
	r.Y.Label.Text = "时间"
	r.NominalX(p.base...)
	err := plotutil.AddLinePoints(r, points...)
	if err != nil {
		log.Println("err when add lines ", err)
		return err
	}
	name := fmt.Sprintf("result_graph %d.png", time.Now().UnixNano())
	if err = r.Save(vg.Length(p.size[0])*vg.Inch, vg.Length(p.size[1])*vg.Inch, name); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
