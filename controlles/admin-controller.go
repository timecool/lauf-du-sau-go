package controlles

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"lauf-du-sau/database"
	"lauf-du-sau/models"
	"lauf-du-sau/service"
	"net/http"
	"time"
)

func ActivateUser(c *gin.Context) {

	userUuid := c.Param("uuid")

	userConnection := database.InitUserCollection()

	update := bson.D{{"$set", bson.D{{"role", models.RoleMember}}}}

	_, err := userConnection.UpdateByID(database.Ctx, userUuid, update)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully User RunActivate"})
}

func ChangeRunStatus(c *gin.Context) {
	userUuid, err := service.GetUserByToken(c)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	runUuid := c.Param("uuid")

	var runStatus models.RunStatusResponse
	if err := c.ShouldBindJSON(&runStatus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if runStatus.Status > 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status not valid"})
		return
	}

	query := bson.M{
		"runs._id": runUuid,
	}

	update := bson.M{
		"$set": bson.M{
			"runs.$.status": runStatus.Status,
		},
	}

	userConnection := database.InitUserCollection()

	result, err := userConnection.UpdateOne(database.Ctx, query, update)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if result.MatchedCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Run was not found"})
		return
	}

	if runStatus.Message != "" {
		push := bson.M{"$push": bson.M{
			"runs.$.messages": models.RunMessage{
				CreateUserUuid: userUuid,
				CreateAt:       time.Now(),
				Message:        runStatus.Message},
		}}
		_, err = userConnection.UpdateOne(database.Ctx, query, push)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Run was successfully edited"})

}

func VerifyRuns(c *gin.Context) {
	userCollection := database.InitUserCollection()

	o1 := bson.M{
		"$unwind": "$runs",
	}
	o2 := bson.M{
		"$match": bson.M{"runs.status": models.RunVerify},
	}
	o3 := bson.M{
		"$group": bson.M{
			"_id":       "$_id",
			"username":  bson.M{"$first": "$username"},
			"email":     bson.M{"$first": "$email"},
			"image_url": bson.M{"$first": "$image_url"},
			"runs":      bson.M{"$push": "$runs"},
		},
	}

	cursor, err := userCollection.Aggregate(database.Ctx, []bson.M{o1, o2, o3})
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

func ResetPassword(c *gin.Context) {
	userUuid := c.Param("uuid")

	var password models.ResetPassword
	if err := c.ShouldBindJSON(&password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userCollection := database.InitUserCollection()
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password.Password), 14)
	update := bson.D{{"$set", bson.M{"password": hashPassword}}}

	result, err := userCollection.UpdateByID(database.Ctx, userUuid, update)
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
