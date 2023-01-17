package controlles

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	user, err := service.GetUserByContext(c)

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	folder := "uploads/" + user.ID.Hex() + "/"
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
	date, err := time.Parse("2006-01-02", dateForm)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	run := models.Run{
		ID:       primitive.NewObjectID(),
		UserID:   user.ID,
		Distance: distance,
		Time:     runTime,
		Date:     date,
		CreateAt: time.Now(),
		Url:      filePath,
		Status:   models.RunVerify,
		Messages: []models.RunMessage{},
	}

	runCollection := database.InitRunCollection()
	_, err = runCollection.InsertOne(database.Ctx, run)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	run.Url = os.Getenv("IMAGE_PATH") + filePath
	c.JSON(http.StatusOK, run)
}

func DeleteRun(c *gin.Context) {

	id := c.Param("uuid")
	permission, run := service.HasPermission(id, c)
	if !permission {
		return
	}
	runCollection := database.InitRunCollection()

	_, err := runCollection.DeleteOne(database.Ctx, bson.M{"_id": run.ID})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Run has been deleted"})

}

func MyRuns(c *gin.Context) {
	user, err := service.GetUserByContext(c)
	month := c.Query("month")
	firstOfMonth, lastOfMonth, err := service.GetFirstAndLastDayFromMonth(month)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	runCollection := database.InitRunCollection()

	query := bson.M{"user_id": user.ID, "date": bson.M{
		"$gte": firstOfMonth,
		"$lte": lastOfMonth,
	}}

	cursor, err := runCollection.Find(database.Ctx, query)
	var runs []models.Run
	if err = cursor.All(database.Ctx, &runs); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	if err := cursor.Close(database.Ctx); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	for index, run := range runs {
		runs[index].Url = os.Getenv("IMAGE_PATH") + run.Url
	}
	c.JSON(http.StatusOK, runs)
}

func TestApi(c *gin.Context) {
	service.OcrImage("https://www.runtastic.com/blog/wp-content/uploads/2018/03/iPhone-Xs-05_ShareImage_de.jpg")
	c.JSON(http.StatusOK, gin.H{"message": "Run check"})

}

func RunsFromUser(c *gin.Context) {
	id := c.Param("uuid")

	month := c.Query("month")
	firstOfMonth, lastOfMonth, err := service.GetFirstAndLastDayFromMonth(month)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	runCollection := database.InitRunCollection()
	userId, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{"user_id": userId, "date": bson.M{
		"$gte": firstOfMonth,
		"$lte": lastOfMonth,
	}}
	cursor, err := runCollection.Find(database.Ctx, query)
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

	id := c.Param("uuid")

	permission, dbRun := service.HasPermission(id, c)
	if !permission {
		return
	}
	var filePath string
	var updateFiles bson.D
	file, header, err := c.Request.FormFile("file")
	if err == nil {
		folder := "uploads/" + dbRun.UserID.String() + "/"
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
		updateFiles = append(updateFiles, bson.E{Key: "url", Value: filePath})
	}

	distance := c.PostForm("distance")
	if distance != "" {
		distanceFloat, err := strconv.ParseFloat(distance, 64)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		updateFiles = append(updateFiles, bson.E{Key: "distance", Value: distanceFloat})
	}

	runTime := c.PostForm("time")
	if runTime != "" {
		timeFloat, err := strconv.ParseFloat(runTime, 64)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		updateFiles = append(updateFiles, bson.E{Key: "time", Value: timeFloat})
	}

	date := c.PostForm("date")
	if date != "" {
		dateForm := c.PostForm("date")
		date, err := time.Parse("2006-01-02", dateForm)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		updateFiles = append(updateFiles, bson.E{Key: "date", Value: date})
	}

	if len(updateFiles) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nothing set"})
		return
	}
	updateFiles = append(updateFiles, bson.E{Key: "status", Value: models.RunVerify})

	update := bson.M{"$set": updateFiles}
	query := bson.M{
		"_id": dbRun.ID,
	}
	runCollection := database.InitRunCollection()

	result, err := runCollection.UpdateOne(database.Ctx, query, update)
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

func NewRuns(c *gin.Context) {

	runCollection := database.InitRunCollection()

	findOptions := options.Find().SetSort(bson.D{{"date", -1}}).SetLimit(4)

	cursor, err := runCollection.Find(database.Ctx, bson.M{}, findOptions)
	var runs []models.Run
	if err = cursor.All(database.Ctx, &runs); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	if err := cursor.Close(database.Ctx); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	var results []models.RunResponse

	for _, run := range runs {
		results = append(results, service.FormatRun(run))
	}

	c.JSON(http.StatusOK, results)
}
