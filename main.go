package main

import (
	"GoLangLoginBackEnd/controllers"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/go-playground/validator.v9"
)

type M map[string]interface{}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	r := echo.New()
	r.Use(middleware.Logger())
	r.Use(middleware.Recover())
	r.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{"*"},
	}))
	r.Validator = &CustomValidator{validator: validator.New()}

	r.POST("/login", controllers.UsersLogin)
	r.POST("/register", controllers.UsersRegister)
	r.POST("/relogin", controllers.UsersReLogin)

	e := r.Group("/api")
	e.Use(middleware.JWT([]byte("secret")))
	e.POST("/logincheck", controllers.UsersLoginCheck)
	e.GET("/", controllers.Index)

	r.Logger.Fatal(r.Start(":9000"))
}
