package ocr

import (
	"context"
	"fmt"
	"testing"

	"github.com/ulysseskk/house/app/common/config"
)

func TestOCRClient_Base64(t *testing.T) {
	client := NewOCRClient(&config.OCRConfig{Host: "http://10.184.21.250:8080"})
	result, err := client.GetByBase64(context.Background(), str)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

var str = `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`
