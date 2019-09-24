package controllers

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo"
)

func Index(c echo.Context) error {
	// Get JWT
	jwt_val := c.Get("user").(*jwt.Token)
	claims := jwt_val.Claims.(jwt.MapClaims)
	username := claims["username"].(string)

	return c.JSON(http.StatusOK, gin.H{
		"data": username,
	})
}
