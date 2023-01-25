package controlles

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"lauf-du-sau/database"
	"lauf-du-sau/service"
	"net/http"
)

func UpdateUser(c *gin.Context) {
	user, err := service.GetUserByContext(c)

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	var filePath string
	var updateFiles bson.D
	file, header, err := c.Request.FormFile("file")
	if err == nil {
		filePath, err = service.SaveProfileImage(header, file, user.ID.Hex())
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
	}
	if filePath != "" {
		user.ImageUrl = filePath
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
		user.Email = email
		updateFiles = append(updateFiles, bson.E{Key: "email", Value: email})
	}

	username := c.PostForm("username")
	if username != "" {
		_, isUsernameSet, _ := service.GetUserByUsername(username, userCollection)
		if isUsernameSet {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
			return
		}
		user.Username = username
		updateFiles = append(updateFiles, bson.E{Key: "username", Value: username})
	}
	//	goal := c.PostForm("goal")
	//	if goal != "" {
	//		updateFiles = append(updateFiles, bson.E{Key: "goal", Value: goal})
	//	}

	if len(updateFiles) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nothing set"})
		return
	}
	update := bson.D{{"$set", updateFiles}}
	_, err = userCollection.UpdateByID(database.Ctx, user.ID, update)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	result := service.UserToResultUser(user)
	c.JSON(http.StatusOK, gin.H{"user": result})
}
func DeleteUserImage(c *gin.Context) {
	user, err := service.GetUserByContext(c)

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	userCollection := database.InitUserCollection()
	update := bson.D{{"$set", bson.M{"image_url": ""}}}
	_, err = userCollection.UpdateByID(database.Ctx, user.ID, update)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	user.ImageUrl = ""
	result := service.UserToResultUser(user)
	c.JSON(http.StatusOK, gin.H{"user": result})
}

func Me(c *gin.Context) {
	user, err := service.GetUserByContext(c)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	result := service.UserToResultUser(user)
	c.JSON(http.StatusOK, gin.H{"user": result})
}

func GetUser(c *gin.Context) {
	userUuid := c.Param("uuid")

	user, _, err := service.GetUserById(userUuid)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	result := service.UserToResultUser(user)
	c.JSON(http.StatusOK, gin.H{"user": result})
}
