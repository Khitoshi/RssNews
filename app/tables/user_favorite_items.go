package tables

import (
	"time"

	"github.com/uptrace/bun"
)

type USER_FAVORITE_ITEMS struct { //ユーザーのお気に入りを格納
	bun.BaseModel `bun:"table:user_favorite_items,alias:ui"`

	User_id    int64     `bun:"user_id,pk"`
	Item_id    int64     `bun:"item_id,pk"`
	Created_at time.Time `bun:"created_at,notnull"`
}
