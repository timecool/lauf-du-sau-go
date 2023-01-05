package routes

import (
	"github.com/gin-gonic/gin"
	"lauf-du-sau/controlles"
	"lauf-du-sau/middleware"
)

const baseApiPattern = "/api/v1"
const user = baseApiPattern + "/user"

func Setup(router *gin.Engine) {
	// api/v1/user
	router.POST(user+"/register", controlles.Register)
	router.POST(user+"/login", controlles.Login)
	router.POST(user+"/run", middleware.Member, controlles.CreateRun)
	router.DELETE(user+"/run", middleware.Member, controlles.DeleteRun)

	router.Static("/uploads", "./uploads")

}
