package controlles

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"lauf-du-sau/database"
	"lauf-du-sau/models"
	"lauf-du-sau/service"
	"net/http"
	"os"
)

func Register(c *gin.Context) {
	// Validate user
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userConnection := database.InitUserCollection()
	_, isEmailSet, _ := service.GetUserByEmail(user.Email, userConnection)
	if isEmailSet {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	// hash and salt the password
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	user.Password = string(hashPassword)
	user.UserRole = models.RoleNone
	// create uuid
	user.UUID = uuid.New().String()

	// save user in collection
	if _, err2 := userConnection.InsertOne(database.Ctx, user); err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
		return
	}
	result := service.UserToResultUser(user)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully registered", "user": result})
}

func Login(c *gin.Context) {
	var login models.User
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userConnection := database.InitUserCollection()

	userInDatabase, isEmailSet, err := service.GetUserByEmail(login.Email, userConnection)

	if !isEmailSet {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User or Password wrong"})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	// check password
	if err := bcrypt.CompareHashAndPassword([]byte(userInDatabase.Password), []byte(login.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User or Password wrong"})
		return
	}
	if userInDatabase.UserRole == "none" {
		// user is not activated
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not yet activated"})
		return
	}

	// create a jwt
	token, err := service.CreateToken(userInDatabase)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.SetCookie("token", token, 1000*60*60*2, "/", os.Getenv("FRONTEND_URL"), false, true)
	result := service.UserToResultUser(userInDatabase)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully registered", "user": result})
}
