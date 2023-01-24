package middleware

import (
	"github.com/gin-gonic/gin"
	"lauf-du-sau/models"
	"lauf-du-sau/service"
	"net/http"
)

func Admin(c *gin.Context) {
	tokenString, _ := c.Cookie("token")

	role, err := service.GetCurrentUserRole(tokenString)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	if role != models.RoleAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you are not an admin"})
		c.Abort()
		return
	}
	c.Next()
}

func Member(c *gin.Context) {
	tokenString, _ := c.Cookie("token")

	role, err := service.GetCurrentUserRole(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	if role == models.RoleNone {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you are not a member"})
		c.Abort()
		return
	}
	c.Next()

}
