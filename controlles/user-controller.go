package controlles

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"lauf-du-sau/database"
	"lauf-du-sau/service"
	"net/http"
)

func UpdateUser(c *gin.Context) {
	userUuid, err := service.GetUserByToken(c)

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	var filePath string
	var updateFiles bson.D
	file, header, err := c.Request.FormFile("file")
	if err == nil {
		filePath, err = service.SaveProfileImage(header, file, userUuid)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
	}
	if filePath != "" {
		updateFiles = append(updateFiles, bson.E{Key: "image_url", Value: filePath})
	}
	userCollection := database.InitUserCollection()

	email := c.PostForm("email")
	if email != "" {
		_, isEmailSet, _ := service.GetUserByEmail(email, userCollection)
		if isEmailSet {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
			return
		}
		updateFiles = append(updateFiles, bson.E{Key: "email", Value: email})
	}

	username := c.PostForm("username")
	if username != "" {
		_, isUsernameSet, _ := service.GetUserByUsername(username, userCollection)
		if isUsernameSet {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
			return
		}
		updateFiles = append(updateFiles, bson.E{Key: "username", Value: username})
	}
	goal := c.PostForm("goal")
	if goal != "" {
		updateFiles = append(updateFiles, bson.E{Key: "goal", Value: goal})
	}

	if len(updateFiles) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nothing set"})
		return
	}
	update := bson.D{{"$set", updateFiles}}
	_, err = userCollection.UpdateByID(database.Ctx, userUuid, update)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	user, _, err := service.GetUserById(userUuid, userCollection)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	result := service.UserToResultUser(user)
	c.JSON(http.StatusOK, gin.H{"user": result})
}

func Me(c *gin.Context) {
	userUuid, err := service.GetUserByToken(c)

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	userCollection := database.InitUserCollection()
	user, _, err := service.GetUserById(userUuid, userCollection)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	result := service.UserToResultUser(user)
	c.JSON(http.StatusOK, gin.H{"user": result})
}

func GetUser(c *gin.Context) {
	userUuid := c.Param("uuid")

	userCollection := database.InitUserCollection()
	user, _, err := service.GetUserById(userUuid, userCollection)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	result := service.UserToResultUser(user)
	c.JSON(http.StatusOK, gin.H{"user": result})
}
