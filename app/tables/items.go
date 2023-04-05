package items

import (
	"time"
)

type ITEMS struct { //RSSから入手したアイテムを保管
	//Id           int64
	Url          string
	Title        string
	Description  string
	Author       string
	Published_at time.Time
	Created_at   time.Time
	Updated_at   time.Time
}
