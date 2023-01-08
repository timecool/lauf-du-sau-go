package routes

import (
	"github.com/gin-gonic/gin"
	"lauf-du-sau/controlles"
	"lauf-du-sau/middleware"
)

const baseApiPattern = "/api/v1"
const user = baseApiPattern + "/user"
const statistics = baseApiPattern + "/statistics"

func Setup(router *gin.Engine) {
	// api/v1/user
	router.POST(user+"/register", controlles.Register)
	router.POST(user+"/login", controlles.Login)

	router.GET(user+"/me", middleware.Member, controlles.Me)
	router.GET(user+"/:uuid", middleware.Member, controlles.GetUser)
	router.PATCH(user, middleware.Member, controlles.UpdateUser)

	// api/v1/user/run
	router.GET(user+"/runs", middleware.Member, controlles.MyRuns)
	router.GET(user+"/runs/:uuid", middleware.Member, controlles.RunsFromUser)

	router.POST(user+"/run", middleware.Member, controlles.CreateRun)
	router.DELETE(user+"/run/:uuid", middleware.Member, controlles.DeleteRun)

	// api/v1/statistics
	router.GET(statistics+"/total-run", middleware.Member, controlles.TotalRun)
	router.GET(statistics+"/leaderboard", middleware.Member, controlles.Leaderboard)

	// only admin
	router.PATCH(user+"/run/:uuid/status", middleware.Admin, controlles.ChangeRunStatus)
	router.PATCH(user+"/:uuid/activate", middleware.Admin, controlles.ActivateUser)

	router.GET("/test", controlles.TestApi)
	router.Static("/uploads", "./uploads")
	router.Static("/profile", "./profile")

}
