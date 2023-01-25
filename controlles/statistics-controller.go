package controlles

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"lauf-du-sau/database"
	"lauf-du-sau/models"
	"lauf-du-sau/service"
	"net/http"
)

func Leaderboard(c *gin.Context) {
	month := c.Query("month")
	var userIDS []models.LeaderboardUserID
	var err error
	if month == "" {
		userIDS, err = service.GetAllTimeLeaderboard()

	} else {
		userIDS, err = service.GetMonthLeaderboard(month)
	}
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	var results []models.LeaderboardUser
	for _, element := range userIDS {
		user, _, _ := service.GetUserById(element.UserID)
		resultUser := service.UserToResultUser(user)
		results = append(results, models.LeaderboardUser{User: resultUser, Total: service.ToFixed(element.Total, 2)})
	}

	c.JSON(http.StatusOK, results)
}

func GetTotal(c *gin.Context) {
	user, err := service.GetUserByContext(c)

	var results []models.LeaderboardUserID

	o1 := bson.M{
		"$match": bson.M{"status": models.RunActivate, "user_id": user.ID},
	}

	o2 := bson.M{
		"$group": bson.M{
			"_id":   "$user_id",
			"total": bson.M{"$sum": "$distance"},
		},
	}

	runCollection := database.InitRunCollection()
	cursor, err := runCollection.Aggregate(database.Ctx, []bson.M{o1, o2})

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return

	}

	if err = cursor.All(database.Ctx, &results); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return

	}
	if err := cursor.Close(database.Ctx); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	if len(results) > 0 {
		totalResult := results[0]
		totalResult.Total = service.ToFixed(totalResult.Total, 2)
		c.JSON(http.StatusOK, totalResult)
	}

}

func GetRunsGroupByDay(c *gin.Context) {
	month := c.Query("month")
	user, err := service.GetUserByContext(c)

	var results []models.RunsGroupByDay
	firstOfMonth, lastOfMonth, err := service.GetFirstAndLastDayFromMonth(month)

	o1 := bson.M{
		"$match": bson.M{"user_id": user.ID},
	}
	o2 := bson.M{
		"$match": bson.M{"date": bson.M{
			"$gte": firstOfMonth,
			"$lte": lastOfMonth,
		}},
	}
	o3 := bson.M{
		"$group": bson.M{
			"_id":   "$date",
			"total": bson.M{"$sum": "$distance"},
			"runs":  bson.M{"$push": "$$ROOT"},
		},
	}

	runCollection := database.InitRunCollection()
	cursor, err := runCollection.Aggregate(database.Ctx, []bson.M{o1, o2, o3})

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	if err = cursor.All(database.Ctx, &results); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	if err := cursor.Close(database.Ctx); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
	return

}
