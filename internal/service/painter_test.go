package service

import (
	"testing"
)

func TestPainter_Load(t *testing.T) {
	p := NewPainter(150, 40)
	if err := p.Load("./paint.csv"); err != nil {
		t.Error(err)
	}
	err := p.Paint()
	if err != nil {
		t.Error(err)
	}

}
