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
	"strconv"
	"time"
)

func CreateRun(c *gin.Context) {
	userUuid, err := service.GetUserByToken(c)
	fmt.Println(userUuid)
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

	run := models.Run{UUID: uuid.New().String(), Distance: distance, Time: runTime, Date: time.Unix(timestamp, 0), CreateAt: time.Now(), Url: filePath}

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

func DeleteRun(c *gin.Context) {

	userUuid, err := service.GetUserByToken(c)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	runUuid, isSet := c.GetQuery("uuid")
	if !isSet {
		c.JSON(http.StatusBadGateway, gin.H{"error": "No run UUID set"})
		return
	}

	userCollection := database.InitUserCollection()
	change := bson.M{"$pull": bson.M{"runs": runUuid}}
	_, err = userCollection.UpdateByID(database.Ctx, userUuid, change)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Run has been deleted"})

}
