package framework

import (
	"testing"
	"time"

	"github.com/abyss414/house/app/common/constant"
)

func Test_ZiroomScrapper(t *testing.T) {
	s := NewScraper("", nil, constant.PlatformZiroom)
	s.Start()
	for !s.finished {
		time.Sleep(2 * time.Second)
	}
}

func Test_LianjiaScrapper(t *testing.T) {
	s := NewScraper("", nil, constant.PlatformLianjia)
	s.Start()
	for !s.finished {
		time.Sleep(2 * time.Second)
	}
}

func Test_5i5jScrapper(t *testing.T) {
	s := NewScraper("", nil, constant.Platform5i5j)
	s.Start()
	for !s.finished {
		time.Sleep(2 * time.Second)
	}
}
