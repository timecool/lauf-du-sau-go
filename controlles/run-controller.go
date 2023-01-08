package controlles

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"lauf-du-sau/database"
	"lauf-du-sau/models"
	"lauf-du-sau/service"
	"lauf-du-sau/utils"
	"net/http"
	"os"
	"strconv"
	"time"
)

func CreateRun(c *gin.Context) {
	userUuid, err := service.GetUserByToken(c)

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	folder := "uploads/" + userUuid + "/"
	err = utils.CreateFolder(folder)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	filePath, err := service.SaveImage(header, file, folder)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	distance, err := strconv.ParseFloat(c.PostForm("distance"), 64)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	runTime, err := strconv.ParseFloat(c.PostForm("time"), 64)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	dateForm := c.PostForm("date")
	date, err := time.Parse("2006-01-02 15:04:05", dateForm)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	run := models.Run{
		UUID:     uuid.New().String(),
		Distance: distance,
		Time:     runTime,
		Date:     date,
		CreateAt: time.Now(),
		Url:      os.Getenv("IMAGE_PATH") + filePath,
		Status:   models.RunVerify,
		Messages: []models.RunMessage{},
	}

	userCollection := database.InitUserCollection()
	change := bson.M{"$push": bson.M{"runs": run}}
	_, err = userCollection.UpdateByID(database.Ctx, userUuid, change)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"run": run})
}

func DeleteRun(c *gin.Context) {
	userUuid, err := service.GetUserByToken(c)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	runUuid := c.Param("uuid")

	userCollection := database.InitUserCollection()
	change := bson.M{"$pull": bson.M{"runs": bson.M{"_id": runUuid}}}
	result, err := userCollection.UpdateByID(database.Ctx, userUuid, change)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	if result.MatchedCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Run was not found"})
		return
	}
	if result.ModifiedCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nothing deleted"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Run has been deleted"})

}

func MyRuns(c *gin.Context) {
	userUuid, err := service.GetUserByToken(c)
	month := c.Query("month")
	firstOfMonth, lastOfMonth, err := service.GetFirstAndLastDayFromMonth(month)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	userCollection := database.InitUserCollection()

	o1 := bson.M{
		"$match": bson.M{"_id": userUuid},
	}
	o2 := bson.M{
		"$unwind": "$runs",
	}
	o3 := bson.M{
		"$match": bson.M{"runs.date": bson.M{
			"$gte": firstOfMonth,
			"$lte": lastOfMonth,
		}},
	}
	o4 := bson.M{
		"$group": bson.M{
			"_id":  "$_id",
			"runs": bson.M{"$push": "$runs"},
		},
	}

	cursor, err := userCollection.Aggregate(database.Ctx, []bson.M{o1, o2, o3, o4})
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

func TestApi(c *gin.Context) {
	service.OcrImage("https://www.runtastic.com/blog/wp-content/uploads/2018/03/iPhone-Xs-05_ShareImage_de.jpg")
	c.JSON(http.StatusOK, gin.H{"message": "Run check"})

}

func RunsFromUser(c *gin.Context) {
	userUuid := c.Param("uuid")

	month := c.Query("month")
	firstOfMonth, lastOfMonth, err := service.GetFirstAndLastDayFromMonth(month)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	userCollection := database.InitUserCollection()

	o1 := bson.M{
		"$match": bson.M{"_id": userUuid},
	}
	o2 := bson.M{
		"$unwind": "$runs",
	}
	o3 := bson.M{
		"$match": bson.M{"runs.status": models.RunActivate},
	}
	o4 := bson.M{
		"$match": bson.M{"runs.date": bson.M{
			"$gte": firstOfMonth,
			"$lte": lastOfMonth,
		}},
	}
	o5 := bson.M{
		"$group": bson.M{
			"_id":  "$_id",
			"runs": bson.M{"$push": "$runs"},
		},
	}

	cursor, err := userCollection.Aggregate(database.Ctx, []bson.M{o1, o2, o3, o4, o5})
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

func AllRuns(c *gin.Context) {
	month := c.Query("month")
	firstOfMonth, lastOfMonth, err := service.GetFirstAndLastDayFromMonth(month)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	userCollection := database.InitUserCollection()

	o1 := bson.M{
		"$unwind": "$runs",
	}
	o2 := bson.M{
		"$match": bson.M{"runs.status": models.RunActivate},
	}
	o3 := bson.M{
		"$match": bson.M{"runs.date": bson.M{
			"$gte": firstOfMonth,
			"$lte": lastOfMonth,
		}},
	}
	o4 := bson.M{
		"$group": bson.M{
			"_id":       "$_id",
			"username":  bson.M{"$first": "$username"},
			"email":     bson.M{"$first": "$email"},
			"image_url": bson.M{"$first": "$image_url"},
			"runs":      bson.M{"$push": "$runs"},
		},
	}

	cursor, err := userCollection.Aggregate(database.Ctx, []bson.M{o1, o2, o3, o4})
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

func UpdateRun(c *gin.Context) {
	tokenString, _ := c.Cookie("token")

	runUuid := c.Param("uuid")

	userRole, err := service.GetCurrentUserRole(tokenString)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	userUuid, err := service.GetCurrentUserUuid(tokenString)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	runUserUuid, err := service.GetUserFromRunUuid(runUuid)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	if runUserUuid == "" {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Run not found"})
		return
	}

	if userRole == models.RoleMember && userUuid != runUserUuid {
		c.JSON(http.StatusForbidden, gin.H{"error": "you have no rights to edit this run"})
		return
	}
	var filePath string
	var updateFiles bson.D
	file, header, err := c.Request.FormFile("file")
	if err == nil {
		folder := "uploads/" + runUserUuid + "/"
		err = utils.CreateFolder(folder)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		filePath, err = service.SaveImage(header, file, folder)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
	}
	if filePath != "" {
		updateFiles = append(updateFiles, bson.E{Key: "runs.$.url", Value: filePath})
	}

	distance := c.PostForm("distance")
	if distance != "" {
		distanceFloat, err := strconv.ParseFloat(distance, 64)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		updateFiles = append(updateFiles, bson.E{Key: "runs.$.distance", Value: distanceFloat})
	}

	runTime := c.PostForm("time")
	if runTime != "" {
		timeFloat, err := strconv.ParseFloat(runTime, 64)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		updateFiles = append(updateFiles, bson.E{Key: "runs.$.time", Value: timeFloat})
	}

	date := c.PostForm("date")
	if date != "" {
		dateForm := c.PostForm("date")
		date, err := time.Parse("2006-01-02 15:04:05", dateForm)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		updateFiles = append(updateFiles, bson.E{Key: "runs.$.date", Value: date})
	}

	if len(updateFiles) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nothing set"})
		return
	}
	updateFiles = append(updateFiles, bson.E{Key: "runs.$.status", Value: models.RunVerify})

	update := bson.M{"$set": updateFiles}

	query := bson.M{
		"runs._id": runUuid,
	}
	userCollection := database.InitUserCollection()

	result, err := userCollection.UpdateOne(database.Ctx, query, update)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if result.MatchedCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Run was not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Run was successfully edited"})

}
