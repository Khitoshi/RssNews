package pageHandles

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"rss_reader/modules"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

func HandleLogin_Get(c echo.Context) error {
	return c.Render(http.StatusOK, "login", nil)
}

func HandleLogin_Post(c echo.Context) error {
	userparam := &modules.User{}
	userparam.Email = c.FormValue("mail")
	userparam.Password = c.FormValue("password")

	//TODO ログイン時の取得情報をidだけに変更する
	u, err := loginUser(userparam)
	if err != nil {
		c.Logger().Fatal(err)
	}

	//クッキーを設定
	cookie := new(http.Cookie)
	cookie.Name = "userId"
	cookie.Value = strconv.FormatInt(u.Id, 10)
	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, "login")
}

type USER struct { //ユーザー情報
	bun.BaseModel `bun:"table:users,alias:u"`

	Id         int64     `bun:"id,pk,autoincrement"`
	Name       string    `bun:"name,notnull"`
	Email      string    `bun:"email,notnull,unique"`
	Password   string    `bun:"password,notnull"`
	Created_at time.Time `bun:"created_at,notnull"`
	Updated_at time.Time `bun:"updated_at,notnull"`
}

// userテーブルで成否check
func loginUser(user *modules.User) (USER, error) {
	//dbを開く
	sqldb, err := sql.Open("postgres", "user=postgres dbname=rss_reader_web password=985632 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer sqldb.Close()

	db := bun.NewDB(sqldb, pgdialect.New())
	defer db.Close()

	//クエリのパラメーター出力
	db.AddQueryHook(bundebug.NewQueryHook(
		//bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	u := USER{}
	err = db.NewSelect().Model(&u).Where("email=? and password=?", user.Email, user.Password).Scan(context.Background())
	if err != nil && err != sql.ErrNoRows {
		return USER{}, err
	}

	return u, nil
}
