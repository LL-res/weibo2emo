package main

import (
	"flag"
	"github.com/emomo/weibo2emo/internal/service"
	"log"
)

var (
	PostsPath      string
	DictionaryPath string
	ResultsPath    string
	ResultTimePath string
)

func main() {
	flag.StringVar(&PostsPath, "postpath", "", "博文数据(.csv)所在的路径 示例 ./原始博文1月.csv")
	flag.StringVar(&DictionaryPath, "dictpath", "", "字典集数据(.csv)所在的路径 示例 ./LIWC.csv")
	flag.StringVar(&ResultsPath, "respath", "./result.csv", "结果数据（未按时间聚合）想要导出的目录，默认导出在当前目录，名称为result.csv,如想自定义导出，一定要精确到文件名")
	flag.StringVar(&ResultTimePath, "resTpath", "./result_time.csv", "按时间聚合的结果数据集想要导出的目录，默认导出在当前目录，名称为result_time.csv，如想自定义导出，一定要精确到文件名")
	flag.Parse()

	p := service.NewProcessor()
	defer p.Slicer.Free()

	err := p.LoadPosts(PostsPath)
	if err != nil {
		log.Println("导入博文集时出现错误 : ", err)
		return
	}

	err = p.LoadDictionary(DictionaryPath)
	if err != nil {
		log.Println("导入字典集时出现错误 : ", err)
		return
	}

	err = p.ExportResult(ResultsPath)
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
