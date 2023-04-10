package pageHandles

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleLogout_Get(c echo.Context) error {
	return c.Render(http.StatusOK, "logout", nil)
}

func HandleLogout_Post(c echo.Context) error {
	cmd := c.FormValue("command")

	if cmd == "del cookie" {
		cookie := new(http.Cookie)
		cookie.Name = "userId"
		cookie.MaxAge = -1
		c.SetCookie(cookie)
	}

	return c.Redirect(http.StatusFound, "login")
}
