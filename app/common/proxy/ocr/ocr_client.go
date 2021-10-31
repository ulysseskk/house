package ocr

import (
	"context"
	"encoding/json"

	"github.com/abyss414/house/app/common/config"
	"github.com/go-resty/resty/v2"
)

type OCRClient struct {
	restyClient *resty.Client
}

func NewOCRClient(conf *config.OCRConfig) *OCRClient {
	restyClient := resty.New()
	restyClient.SetHostURL(conf.Host).SetHeader("content-type", "application/json")
	return &OCRClient{restyClient}
}

func (client *OCRClient) GetByBase64(ctx context.Context, fileBase64 string) (string, error) {
	resp, err := client.restyClient.R().SetContext(ctx).SetBody(&Base64Request{
		Base64: fileBase64,
		Trim:   "\n",
	}).Post("/base64")
	if err != nil {
		return "", err
	}
	rawResponse := resp.String()
	jsonResult := &FileResponse{}
	if err := json.Unmarshal([]byte(rawResponse), jsonResult); err != nil {
		return "", err
	}
	return jsonResult.Result, nil
}
