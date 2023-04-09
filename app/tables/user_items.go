package tables

import (
	"time"

	"github.com/uptrace/bun"
)

type USER_ITEMS struct { //ユーザーとユーザーが登録した記事をつなぐ
	bun.BaseModel `bun:"table:user_items,alias:ui"`

	User_id    int64     `bun:"user_id,pk"`
	Rss_id     int64     `bun:"rss_id,pk"`
	Created_at time.Time `bun:"created_at,notnull"`
}
