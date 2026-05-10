package util

import (
	"kenaito-dns/config"
	"time"
)

func NowStr() string {
	return time.Now().Format(config.AppTimeFormat)
}
