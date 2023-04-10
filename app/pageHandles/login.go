package pageHandles

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"rss_reader/tables"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

func HandleLogin_Get(c echo.Context) error {
	return c.Render(http.StatusOK, "login", nil)
}

func HandleLogin_Post(c echo.Context) error {
	userparam := &tables.USER{}
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

	return c.Redirect(http.StatusFound, "home")
}

// userテーブルで成否check
func loginUser(user *tables.USER) (tables.USER, error) {
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

	u := tables.USER{}
	err = db.NewSelect().Model(&u).Where("email=? and password=?", user.Email, user.Password).Scan(context.Background())
	if err != nil && err != sql.ErrNoRows {
		return tables.USER{}, err
	}

	return u, nil
}
