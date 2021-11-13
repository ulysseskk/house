package lianjia

import (
	"context"
	"fmt"
	"testing"

	"github.com/ulysseskk/house/app/common/model/db"
)

func TestScrapper_GetAllAreaLink(t *testing.T) {
	s := NewLianjiaExecutor()
	areas, err := s.GetAllAreaLink(context.Background(), "bj")
	if err != nil {
		panic(err)
	}
	for area, url := range areas {
		fmt.Println(area, url)
	}
}

func TestScrapper_GetSinglePage(t *testing.T) {
	s := NewLianjiaExecutor()
	nextPage, err := s.executeForPage(context.Background(), "https://bj.lianjia.com/ershoufang/", []executeFunc{func(ctx context.Context, detail *db.ErShouFangDetail) error {
		fmt.Println(detail)
		return nil
	}})
	if err != nil {
		panic(err)
	}
	fmt.Println(nextPage)
}
