package pageHandles

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"rss_reader/modules"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

func HandleSignup_Get(c echo.Context) error {
	return c.Render(http.StatusOK, "signup", nil)
}

// ユーザー登録するページ
func HandleSignup_Post(c echo.Context) error {
	userparam := &modules.SignUpInput{}
	userparam.Name = c.FormValue("name")
	userparam.Email = c.FormValue("mail")
	userparam.Password = c.FormValue("password")

	err := registrationUser(userparam)
	if err != nil {
		c.Logger().Fatal(err)
	}

	return c.Redirect(http.StatusFound, "login")
}

// signup情報をデータベースに登録
func registrationUser(userInfo *modules.SignUpInput) error {
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

	//TODO:パスワードハッシュ化
	u := modules.User{
		Name:       userInfo.Name,
		Email:      userInfo.Email,
		Password:   userInfo.Password,
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}

	_, err = db.NewInsert().Model(&u).Exec(context.Background())
	if err != nil {
		return err
	}

	fmt.Printf("insert %v\n", u.Name)
	return nil
}
