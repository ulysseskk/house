package db

type SystemCache struct {
	CacheName  string `db:"cache_name" json:"cache_name"`
	CacheValue string `db:"cache_value" json:"cache_value"`
}

func (s SystemCache) TableName() string {
	return "system_cache"
}
