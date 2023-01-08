package controlles

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"lauf-du-sau/database"
	"lauf-du-sau/models"
	"lauf-du-sau/service"
	"net/http"
	"os"
)

func Leaderboard(c *gin.Context) {
	month := c.Query("month")
	userCollection := database.InitUserCollection()
	var results []models.LeaderboardUser
	var err error
	if month == "" {
		results, err = service.GetAllTimeLeaderboard(userCollection)

	} else {
		results, err = service.GetMonthLeaderboard(userCollection, month)
	}
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	for index, element := range results {
		if element.ImageUrl != "" {
			results[index].ImageUrl = os.Getenv("IMAGE_PATH") + element.ImageUrl

		}
	}

	c.JSON(http.StatusOK, results)
}

func TotalRun(c *gin.Context) {
	userCollection := database.InitUserCollection()
	o1 := bson.M{
		"$unwind": "$runs",
	}
	o2 := bson.M{
		"$group": bson.M{
			"_id":   "null",
			"total": bson.M{"$sum": "$runs.distance"},
		},
	}

	cursor, err := userCollection.Aggregate(database.Ctx, []bson.M{o1, o2})

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	var results []bson.M
	if err = cursor.All(database.Ctx, &results); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	if err := cursor.Close(database.Ctx); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results[0])

}
