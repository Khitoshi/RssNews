package pageHandles

import (
	"context"
	"fmt"
	"net/http"
	"rss_reader/cipher"
	"rss_reader/database"
	"rss_reader/tables"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

// ユーザーが登録画面で設定するパラメーター
type SignUpInput struct {
	Name     string
	Email    string
	Password string
}

func HandleSignup_Get(c echo.Context) error {
	return c.Render(http.StatusOK, "signup", nil)
}

// ユーザー登録するページ
func HandleSignup_Post(c echo.Context) error {
	userparam := &SignUpInput{}
	userparam.Name = c.FormValue("name")
	userparam.Email = c.FormValue("mail")
	userparam.Password = cipher.HashStr(c.FormValue("password"))

	err := registrationUser(userparam)
	if err != nil {
		c.Logger().Fatal(err)
	}

	return c.Redirect(http.StatusFound, "login")
}

// signup情報をデータベースに登録
func registrationUser(userInfo *SignUpInput) error {

	u := tables.USER{
		Name:       userInfo.Name,
		Email:      userInfo.Email,
		Password:   userInfo.Password,
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}

	//userテーブルに登録 登録
	err := database.WithDBConnection(func(db *bun.DB) error {
		_, err := db.NewInsert().Model(&u).Exec(context.Background())
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	fmt.Printf("insert %v\n", u.Name)
	return nil
}
