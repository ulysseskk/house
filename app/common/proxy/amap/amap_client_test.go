package amap

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
)

func TestAmapClient_SearchPOI(t *testing.T) {
	client := NewAmapClient("https://restapi.amap.com/", "4d568b8b4180e26a960b0a1968268755")
	resp, err := client.SearchPOI(context.Background(), []string{"玺源台"}, "", "北京", true, 1, 10, "")
	if err != nil {
		panic(err)
	}
	jsonByte, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonByte))
}
