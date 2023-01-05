package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"lauf-du-sau/models"
	"lauf-du-sau/service"
	"net/http"
)

func Admin(c *gin.Context) {

}

func Member(c *gin.Context) {
	tokenString, _ := c.Cookie("token")

	fmt.Println("token:" + tokenString)
	role, err := service.GetCurrentUserRole(tokenString)
	if err != nil || role == models.RoleNone {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Next()

}
