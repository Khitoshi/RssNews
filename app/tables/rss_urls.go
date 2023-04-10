package tables

import (
	"time"

	"github.com/uptrace/bun"
)

type RSS_URLS struct {
	bun.BaseModel `bun:"table:rss_urls,alias:ru"`
	//User_id    int64     `bun:"user_id,pk"`

	Rss_id     int64     `bun:"rss_id,pk,autoincrement"`
	Rss_URL    string    `bun:"rss_url,notnull"`
	Created_at time.Time `bun:"created_at,notnull"`
}
