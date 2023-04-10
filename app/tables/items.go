package tables

import (
	"time"

	"github.com/uptrace/bun"
)

type ITEMS struct { //RSSから入手したアイテムを保管
	bun.BaseModel `bun:"table:items,alias:i"`

	Id           int64     `bun:"id,pk,autoincrement"`
	Rss_id       int64     `bun:"rss_id,notnull"`
	Url          string    `bun:"url,notnull,unique"`
	Title        string    `bun:"title,notnull"`
	Description  string    `bun:"description,notnull"`
	Author       string    `bun:"author"`
	Published_at time.Time `bun:"published_at"`
	Created_at   time.Time `bun:"created_at,notnull"`
	Updated_at   time.Time `bun:"updated_at,notnull"`
}
