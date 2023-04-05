package tables

import (
	"time"
)

type RSS_URLS struct {
	Rss_id     int64
	Rss_URL    string
	Created_at time.Time
}
