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
	err = p.LoadDictionary("C:\\Users\\LL\\Desktop\\LIWC.csv")
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
