package amap

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/abyss414/house/app/common/log"

	"github.com/go-resty/resty/v2"
)

func NewAmapClient(host string, key string) *AmapClient {
	restyClient := resty.New().SetHostURL(host)
	restyClient.OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {

		ctx := request.Context()
		values := map[string]interface{}{}
		bodyByte, err := json.Marshal(request.Body)
		if err != nil {
			values["body"] = request.Body
		} else {
			values["body"] = string(bodyByte)
		}
		values["header"] = request.Header.Clone()
		values["url"] = request.URL
		values["method"] = request.Method
		values["query"] = request.QueryParam
		log.WithContext(ctx).WithFields(values).Infof("Request %s", "amap")
		return nil
	})
	restyClient.OnAfterResponse(func(client *resty.Client, response *resty.Response) error {
		ctx := response.Request.Context()
		values := map[string]interface{}{}
		values["header"] = response.Header().Clone()
		values["body"] = response.String()
		values["status_code"] = response.StatusCode()
		values["status"] = response.Status()
		log.WithContext(ctx).WithFields(values).Infof("Request %s receive response", "amap")
		return nil
	})
	return &AmapClient{
		restyClient: restyClient,
		key:         key,
	}
}

type AmapClient struct {
	restyClient *resty.Client
	key         string
}

func (client *AmapClient) SearchPOI(ctx context.Context,
	keywords []string,
	types string,
	city string,
	cityLimit bool,
	pageNum,
	pageSize int,
	extensions string) (*POIResponse, error) {
	queryMap := map[string]string{
		"key": client.key,
	}
	if len(keywords) != 0 {
		queryMap["keywords"] = strings.Join(keywords, "|")
	}
	if len(types) != 0 {
		queryMap["types"] = types
	}
	if city == "" {
		queryMap["city"] = city
	}
	queryMap["citylimit"] = strconv.FormatBool(cityLimit)
	queryMap["offset"] = strconv.Itoa(pageSize)
	queryMap["page"] = strconv.Itoa(pageNum)
	if extensions != "" {
		queryMap["extensions"] = extensions
	}
	resp, err := client.restyClient.R().SetContext(ctx).SetResult(&POIResponse{}).SetQueryParams(queryMap).Get("v3/place/text")
	if err != nil {
		return nil, err
	}
	rawResponse := resp.Result().(*POIResponse)
	return rawResponse, nil
}
