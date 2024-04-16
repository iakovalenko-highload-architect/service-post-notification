package middlewares

import (
	"strings"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

func (m *Middlewares) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		vals := strings.Fields(c.Request().Header.Get("Authorization"))
		if len(vals) != 2 {
			return websocket.ErrBadHandshake
		}

		userID, err := m.tokenManager.ExtractUserID(vals[1])
		if err != nil {
			return websocket.ErrBadHandshake
		}

		c.Set("userID", userID)
		return next(c)
	}
}
