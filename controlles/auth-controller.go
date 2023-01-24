package controlles

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"lauf-du-sau/database"
	"lauf-du-sau/models"
	"lauf-du-sau/service"
	"net/http"
	"strings"
)

func Register(c *gin.Context) {
	// Validate user
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Email == "" || user.Username == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please fill all fields"})
		return
	}

	userConnection := database.InitUserCollection()
	checkEmail := strings.HasSuffix(user.Email, "@byte5.de")
	if !checkEmail {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please use your work email"})
		return
	}
	_, isEmailSet, _ := service.GetUserByEmail(user.Email, userConnection)
	if isEmailSet {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}
	_, isUserSet, _ := service.GetUserByUsername(user.Username, userConnection)
	if isUserSet {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	// hash and salt the password
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	user.Password = string(hashPassword)
	user.UserRole = models.RoleMember
	user.ID = primitive.NewObjectID()

	// save user in collection
	_, err := userConnection.InsertOne(database.Ctx, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// create a jwt
	token, err := service.CreateToken(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	cookie := service.CreateCookie(token, 1000*60*60*2)
	http.SetCookie(c.Writer, &cookie)
	c.Redirect(http.StatusMovedPermanently, "/")

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

	cookie := service.CreateCookie(token, 1000*60*60*2)
	http.SetCookie(c.Writer, &cookie)
	c.Redirect(http.StatusMovedPermanently, "/")

	result := service.UserToResultUser(userInDatabase)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully registered", "user": result})
}

func Logout(c *gin.Context) {
	cookie := service.CreateCookie("", -1)
	http.SetCookie(c.Writer, &cookie)
	c.Redirect(http.StatusMovedPermanently, "/")

	c.JSON(http.StatusOK, gin.H{"message": "Successfully Logout"})
}
