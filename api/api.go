package api
import (
	"github.com/labstack/echo"
	"net/http"
	mw "github.com/labstack/echo/middleware"
)

func initMiddleware(e *echo.Echo) {
	e.Use(mw.Logger())
	e.Use(mw.Recover())
	e.Use(func (c *echo.Context) error {
		h := c.Request().Header
		ctHeader := h.Get("Content-Type")
		if ctHeader == "" {
			h.Set("Content-Type", "application/json")
		}

		return nil
	})
}

func initRoutes(e *echo.Echo) {
	e.Get("/", func (c *echo.Context) error {
		c.String(http.StatusOK, "haha")
		return nil
	})

	e.Put("/auth", RegisterUserHandler)
}

func Initialize(e *echo.Echo) {
	initMiddleware(e)
	initRoutes(e)
}