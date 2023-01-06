package controlles

import (
	"fmt"
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
	timestamp, err := strconv.ParseInt(c.PostForm("date"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	run := models.Run{
		UUID:     uuid.New().String(),
		Distance: distance,
		Time:     runTime,
		Date:     time.Unix(timestamp, 0),
		CreateAt: time.Now(),
		Url:      os.Getenv("IMAGE_PATH") + filePath,
		Status:   models.Verify,
	}

	userCollection := database.InitUserCollection()
	fmt.Println(run)
	change := bson.M{"$push": bson.M{"runs": run}}
	_, err = userCollection.UpdateByID(database.Ctx, userUuid, change)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"run": run})
}

func ChangeRunStatus(c *gin.Context) {
	runUuid := c.Param("uuid")
	status := c.Query("status")
	statusInt, err := strconv.Atoi(status)
	if err != nil || statusInt > 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status not valid"})
		return
	}

	query := bson.M{
		"runs._id": runUuid,
	}

	update := bson.M{
		"$set": bson.M{
			"runs.$.status": status,
		},
	}
	userConnection := database.InitUserCollection()

	userConnection.FindOneAndUpdate(database.Ctx, query, update)

	c.JSON(http.StatusOK, gin.H{"message": "Successfully Update Run"})
}

func DeleteRun(c *gin.Context) {
	userUuid, err := service.GetUserByToken(c)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	runUuid := c.Param("uuid")

	userCollection := database.InitUserCollection()
	change := bson.M{"$pull": bson.M{"runs": runUuid}}
	_, err = userCollection.UpdateByID(database.Ctx, userUuid, change)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Run has been deleted"})

}

func MyRuns(c *gin.Context) {
	userUuid, err := service.GetUserByToken(c)
	month := c.Query("month")
	date, err := time.Parse("2006-01", month)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	currentYear, currentMonth, _ := date.Date()
	currentLocation := date.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
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
