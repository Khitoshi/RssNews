package table_user_items

import "time"

type USER_ITEMS struct { //ユーザーが登録した，記事情報を保管
	User_id    int64
	Item_id    int64
	Created_at time.Time
}
