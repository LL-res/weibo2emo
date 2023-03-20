package service

import (
	"testing"
)

func TestProcessor_ExportResult(t *testing.T) {
	p := NewProcessor()
	err := p.LoadPosts("C:\\Users\\LL\\Desktop\\原始博文1月.csv")
	if err != nil {
		t.Error(err)
	}
	err = p.LoadDictionary("C:\\Users\\LL\\Desktop\\LIWC1.csv")
	if err != nil {
		t.Error(err)
	}
	err = p.ExportResult("./result.csv")
	if err != nil {
		t.Error(err)
	}
	err = p.ExportResultByTime("./result_time.csv")
	if err != nil {
		t.Error(err)
	}

}
func Test2(t *testing.T) {
	p := NewProcessor()
	p.ToDB = true
	err := p.FillTypeMap("C:/Users/LL/Desktop/博主名称及id.csv")
	if err != nil {
		t.Error(err)
	}
	err = p.LoadDictionary("C:/Users/LL/Desktop/LIWC2.csv")
	if err != nil {
		t.Error(err)
	}
	err = p.LoadPosts("C:/Users/LL/Desktop/原始博文_demo.xlsx")
	if err != nil {
		t.Error(err)
	}
	err = p.ExportResult("./result.csv")
	if err != nil {
		t.Error(err)
	}
	err = p.ExportResultByTime("./result_time.csv")
	if err != nil {
		t.Error(err)
	}
}
