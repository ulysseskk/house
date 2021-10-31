package connector

import (
	"github.com/abyss414/house/app/common/config"
	"github.com/abyss414/house/app/common/proxy/ocr"
)

var OcrConnector *ocr.OCRClient

func InitOcrConnector() {
	OcrConnector = ocr.NewOCRClient(config.GlobalConfig().OCR)
}
