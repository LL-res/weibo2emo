package main

import (
	"flag"
	"log"

	"github.com/emomo/weibo2emo/internal/service"
	"github.com/emomo/weibo2emo/internal/tools"
)

var (
	PostsPath       string
	DictionaryPath  string
	ResultsPath     string
	ResultTimePath  string
	GradePath       string
	ResultGradePath string
	TypePath        string
	p               *service.Processor
	g               *service.GradeParser
	TransFunc       string
	Grade           int

	PrintPath string
	Width     int
	Heigt     int
)

func main() {

	flag.StringVar(&PostsPath, "pp", "", "博文数据(.csv)所在的路径 示例 ./原始博文1月.csv (如要进行词频统计，则必填，否则，非必填)")
	flag.StringVar(&DictionaryPath, "pd", "", "字典集数据(.csv)所在的路径 示例 ./LIWC.csv (如要进行词频统计，则必填，否则，非必填)")
	flag.StringVar(&ResultsPath, "pr", "./result.csv", "结果数据（未按时间聚合）想要导出的目录，默认导出在当前目录，名称为result.csv,如想自定义导出，一定要精确到文件名 (选填)")
	flag.StringVar(&ResultTimePath, "prt", "./result_time.csv", "按时间聚合的结果数据集想要导出的目录，默认导出在当前目录，名称为result_time.csv，如想自定义导出，一定要精确到文件名 (选填)")
	flag.StringVar(&TypePath, "pt", "", "id,用户名数据集(.csv)所在路径 示例 ./博主名称及id.csv (如要进行词频统计，则必填，否则，非必填)")
	flag.StringVar(&GradePath, "gp", "", "待分级的数据集(.csv)所在路径 示例 ./result_time.csv (如要进行分级统计，则必填，否则，非必填)")
	flag.StringVar(&ResultGradePath, "gr", "./result_grade.csv", "获取到的分级结果，默认导出在当前目录，名称为result_grade.csv，如想自定义导出，一定要精确到文件名 (选填)")
	flag.StringVar(&TransFunc, "gf", "linear", "分级时所使用的映射函数，目前可选的有 \n linear : f(x) = x \n log : f(x) = ln(x + 1) \n sigmoid : f(x) = 1 / (1 + exp(-x)) + 0.5 \n tanh : 2 / (1 + exp(-2x)) -1 \n 默认使用线性函数 (选填)")
	flag.IntVar(&Grade, "gn", 5, "所要分出的等级数量，默认为 5 (选填)")

	flag.StringVar(&PrintPath, "rp", "", "待输出绘制的csv文件所在位置 实例./result_time.csv (如要进行折线绘制则必填，否则，非必填)")
	flag.IntVar(&Width, "rw", 150, "待输出的图片长度，默认为150，非必填")
	flag.IntVar(&Heigt, "rh", 40, "待输出的图片高度，默认为40，非必填")
	flag.Parse()

	if GradePath != "" {
		g = service.NewGradeParser(Grade)
		err := g.LoadGradeData(GradePath)
		if err != nil {
			log.Println("导入数据集时出现错误 : ", err)
			return
		}
		err = g.ExportResult(ResultGradePath, tools.F(TransFunc))
		if err != nil {
			log.Println("导出分级结果时出现错误 : ", err)
			return
		}
	}

	if DictionaryPath != "" && PostsPath != "" && TypePath != "" {
		p = service.NewProcessor()
		defer p.Slicer.Free()
	}
	if TypePath != "" {
		err := p.FillTypeMap(TypePath)
		if err != nil {
			log.Println("导入类型集时出现错误 : ", err)
			return
		}
	}
	if PostsPath != "" {
		err := p.LoadPosts(PostsPath)
		if err != nil {
			log.Println("导入博文集时出现错误 : ", err)
			return
		}
	}

	if DictionaryPath != "" {
		err := p.LoadDictionary(DictionaryPath)
		if err != nil {
			log.Println("导入字典集时出现错误 : ", err)
			return
		}
	}

	if DictionaryPath != "" && PostsPath != "" && TypePath != "" {
		err := p.ExportResult(ResultsPath)
		if err != nil {
			log.Println("导出结果集时出现错误 : ", err)
			return
		}
		err = p.ExportResultByTime(ResultTimePath)
		if err != nil {
			log.Println("导出时间结果集时出现错误 : ", err)
			return
		}
	}
	if PrintPath != "" {
		r := service.NewPainter(Width, Heigt)
		err := r.Load(PrintPath)
		if err != nil {
			log.Println("导入csv文件错误 : ", err)
			return
		}
		err = r.Paint()
		if err != nil {
			log.Println("导出图片错误 : ", err)
			return
		}
	}

}
