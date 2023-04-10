package tables

import (
	"time"

	"github.com/uptrace/bun"
)

type USER struct { //ユーザー情報
	bun.BaseModel `bun:"table:users,alias:u"`

	Id         int64     `bun:"id,pk,autoincrement"`
	Name       string    `bun:"name,notnull"`
	Email      string    `bun:"email,notnull,unique"`
	Password   string    `bun:"password,notnull"`
	Created_at time.Time `bun:"created_at,notnull"`
	Updated_at time.Time `bun:"updated_at,notnull"`
}
