package config

import (
	"testing"
)

func TestInitConfig(t *testing.T) {
	err := InitConfig()
	if err != nil {
		panic(err)
	}
}
