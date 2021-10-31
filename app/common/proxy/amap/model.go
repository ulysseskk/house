package amap

type POIResponse struct {
	Status     string                 `json:"status"`
	Info       string                 `json:"info"`
	Count      string                 `json:"count"`
	Suggestion *POIResponseSuggestion `json:"suggestion"`
	Pois       []*Poi                 `json:"pois"`
}

type POIResponseSuggestion struct {
	Keywords []interface{} `json:"keywords"`
	Cities   []*SuggestionCity
}

type SuggestionCity struct {
	Name     string `json:"name"`
	Num      int    `json:"num"`
	CityCode int    `json:"city_codde"`
	ADCode   int    `json:"ad_code"`
}

type Poi struct {
	Id         string                 `json:"id"`
	Parent     interface{}            `json:"parent"`
	Address    interface{}            `json:"address"`
	PName      string                 `json:"pname"`
	Importance interface{}            `json:"importance"`
	BizExt     map[string]interface{} `json:"biz_ext"`
	BizType    interface{}            `json:"biz_type"`
	CityName   string                 `json:"cityname"`
	Type       string                 `json:"type"`
	Photos     []*POIPhoto            `json:"photos"`
	TypeCode   string                 `json:"typecode"`
	ShopInfo   string                 `json:"shop_info"`
	PoiWeight  interface{}            `json:"poi_weight"`
	ChildType  interface{}            `json:"child_type"`
	AdName     string                 `json:"adname"`
	Name       string                 `json:"name"`
	Location   string                 `json:"location"`
	Tel        interface{}            `json:"tel"`
	ShopId     interface{}            `json:"shop_id"`
}

type POIPhoto struct {
	Title interface{} `json:"title"`
	Url   string      `json:"url"`
}
