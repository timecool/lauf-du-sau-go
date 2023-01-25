package routes

import (
	"github.com/gin-gonic/gin"
	"lauf-du-sau/controlles"
	"lauf-du-sau/middleware"
)

const baseApiPattern = "/api/v1"
const user = baseApiPattern + "/user"
const runs = baseApiPattern + "/runs"
const admin = baseApiPattern + "/admin"
const statistics = baseApiPattern + "/statistics"

func Setup(router *gin.Engine) {
	// api/v1/user
	router.POST(user+"/register", controlles.Register)
	router.POST(user+"/login", controlles.Login)
	router.DELETE(user+"/logout", middleware.Member, controlles.Logout)

	router.GET(user+"/me", middleware.Member, controlles.Me)
	router.GET(user+"/:uuid", middleware.Member, controlles.GetUser)
	router.PATCH(user, middleware.Member, controlles.UpdateUser)

	router.GET(user+"/runs", middleware.Member, controlles.MyRuns)
	router.GET(user+"/runs/group", middleware.Member, controlles.GetRunsGroupByDay)

	// api/v1/user/run
	router.POST(user+"/run", middleware.Member, controlles.CreateRun)
	router.DELETE(user+"/run/:uuid", middleware.Member, controlles.DeleteRun)
	router.PATCH(user+"/run/:uuid", middleware.Member, controlles.UpdateRun)

	// api/v1/news
	router.GET(baseApiPattern+"/runs/new", middleware.Member, controlles.NewRuns)

	// api/v1/runs
	router.GET(runs+"/:uuid", middleware.Member, controlles.RunsFromUser)
	router.PATCH(runs+"/:uuid", middleware.Member, controlles.RunsFromUser)

	// api/v1/statistics
	router.GET(statistics+"/leaderboard", middleware.Member, controlles.Leaderboard)
	router.GET(statistics+"/total", middleware.Member, controlles.GetTotal)

	// only admin
	router.PATCH(admin+"/run/:uuid/status", middleware.Admin, controlles.ChangeRunStatus)
	router.PATCH(admin+"/:uuid/reset-password", middleware.Admin, controlles.ResetPassword)
	router.GET(admin+"/runs/verify", middleware.Admin, controlles.VerifyRuns)

	router.GET("/test", controlles.TestApi)
	router.Static("/uploads", "./uploads")
	router.Static("/profile", "./profile")

}
