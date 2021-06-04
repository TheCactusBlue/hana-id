package web

import (
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/thecactusblue/hana-id/auth"
	"gorm.io/gorm"
)

func Inject(e *echo.Echo, db *gorm.DB) {
	authService := auth.NewService(db, auth.NewJWTFactory(viper.GetString("SECRET")))

	e.POST("/auth/register", func(c echo.Context) error {
		body := new(struct {
			Email    string `json:"email"`
			Name     string `json:"name"`
			Password string `json:"password"`
		})
		if err := c.Bind(body); err != nil {
			return err
		}
		err := authService.Register(body.Email, body.Name, body.Password)
		if err != nil {
			return err
		}
		return c.JSON(200, echo.Map{"status": "ok"})
	})
	e.POST("/auth/login", func(c echo.Context) error {
		body := new(struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		})
		if err := c.Bind(body); err != nil {
			return err
		}
		pair, err := authService.Login(body.Email, body.Password)
		if err != nil {
			return err
		}

		return c.JSON(200, pair)
	})

	e.POST("/auth/change-password", func(c echo.Context) error {
		body := new(struct {
			Email       string `json:"email"`
			OldPassword string `json:"old_password"`
			NewPassword string `json:"new_password"`
		})
		if err := c.Bind(body); err != nil {
			return err
		}
		err := authService.ChangePassword(body.Email, body.NewPassword, body.OldPassword)
		if err != nil {
			return err
		}

		return c.JSON(200, echo.Map{"status": "ok"})
	})

	e.POST("/auth/refresh", func(c echo.Context) error {
		body := new(struct {
			Refresh string `json:"refresh"`
		})
		if err := c.Bind(body); err != nil {
			return err
		}
		pair, err := authService.Refresh(body.Refresh)
		if err != nil {
			return err
		}

		return c.JSON(200, pair)
	})
}
