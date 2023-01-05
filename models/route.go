package models

import "github.com/gin-gonic/gin"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc gin.HandlerFunc
	Role        Role
}
