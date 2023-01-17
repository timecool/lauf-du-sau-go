package controlles

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"lauf-du-sau/database"
	"lauf-du-sau/models"
	"lauf-du-sau/service"
	"net/http"
	"time"
)

func ChangeRunStatus(c *gin.Context) {
	user, err := service.GetUserByContext(c)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("uuid")

	var runStatus models.RunStatusResponse
	if err := c.ShouldBindJSON(&runStatus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if runStatus.Status > 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status not valid"})
		return
	}
	runId, err := primitive.ObjectIDFromHex(id)
	query := bson.M{
		"_id": runId,
	}

	update := bson.M{
		"$set": bson.M{
			"status": runStatus.Status,
		},
		"$push": bson.M{
			"messages": models.RunMessage{
				CreateUserUuid: user.ID,
				CreateAt:       time.Now(),
				Message:        runStatus.Message},
		},
	}

	runConnection := database.InitRunCollection()

	result, err := runConnection.UpdateOne(database.Ctx, query, update)
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

func VerifyRuns(c *gin.Context) {
	runCollection := database.InitRunCollection()

	cursor, err := runCollection.Find(database.Ctx, bson.M{"status": models.RunVerify})
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

func ResetPassword(c *gin.Context) {
	id := c.Param("uuid")

	userId, err := primitive.ObjectIDFromHex(id)
	var password models.ResetPassword
	if err := c.ShouldBindJSON(&password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userCollection := database.InitUserCollection()
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password.Password), 14)
	update := bson.D{{"$set", bson.M{"password": string(hashPassword)}}}

	result, err := userCollection.UpdateByID(database.Ctx, userId, update)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	if result.MatchedCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User was not found"})
		return
	}
	if result.ModifiedCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nothing updated"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User password was updated"})

}
