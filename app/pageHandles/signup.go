package pageHandles

import (
	"net/http"
	"rss_reader/modules"

	"github.com/labstack/echo/v4"
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

	//TODO:データベースに登録処理を記述する

	return c.Redirect(http.StatusFound, "signup")
}
