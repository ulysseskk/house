package ziroom

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
)

type PriceItem struct {
	Picture  []byte  `json:"picture"`
	OffsetPX float64 `json:"offset_px"`
}

func CalculatePrice(items []*PriceItem) (float64, error) {
	for i, item := range items {
		imageItem, _, err := image.Decode(bytes.NewReader(item.Picture))
		if err != nil {
			return 0, err
		}
		fmt.Println(i, imageItem)
	}
	return 0, nil
}

func DealWithPic() {
	file1, _ := os.Open("img.png") //打开图片1
	var (
		img1 image.Image
		err  error
	)

	if img1, err = png.Decode(file1); err != nil {
		log.Fatal(err)
		return
	}
	newImg := image.NewNRGBA(image.Rect(0, 0, 20, 28))
	draw.Draw(newImg, newImg.Bounds(), img1, img1.Bounds().Min.Sub(image.Pt(int(-1*math.Floor(192.6/28*20)), 0)), draw.Over)
	imgfile, _ := os.Create("need_be_2_1.jpg")
	defer imgfile.Close()
	jpeg.Encode(imgfile, newImg, &jpeg.Options{100})
}
