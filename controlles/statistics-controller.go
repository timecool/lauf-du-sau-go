package controlles

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"lauf-du-sau/database"
	"net/http"
)

func Leaderboard(c *gin.Context) {
	userCollection := database.InitUserCollection()
	o1 := bson.M{
		"$unwind": "$runs",
	}
	o2 := bson.M{
		"$group": bson.M{
			"_id":   "$_id",
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
