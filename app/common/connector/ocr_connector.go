package connector

import (
	"github.com/ulysseskk/house/app/common/config"
	"github.com/ulysseskk/house/app/common/proxy/ocr"
)

var OcrConnector *ocr.OCRClient

func InitOcrConnector() {
	OcrConnector = ocr.NewOCRClient(config.GlobalConfig().OCR)
}
