package service

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"lauf-du-sau/database"
	"lauf-du-sau/models"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

func GetUserByEmail(email string, usersCollection *mongo.Collection) (models.User, bool, error) {
	var user models.User
	err := usersCollection.FindOne(database.Ctx, bson.M{"email": email}).Decode(&user)
	return user, len(user.Username) != 0, err
}
func GetUserByUsername(username string, usersCollection *mongo.Collection) (models.User, bool, error) {
	var user models.User
	err := usersCollection.FindOne(database.Ctx, bson.M{"username": username}).Decode(&user)
	return user, len(user.Username) != 0, err
}

func GetUserById(id string) (models.User, bool, error) {
	var user models.User
	usersCollection := database.InitUserCollection()
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, false, err
	}
	err = usersCollection.FindOne(database.Ctx, bson.M{"_id": objectId}).Decode(&user)
	return user, len(user.Username) != 0, err
}

func UserToResultUser(user models.User) models.ReturnUser {
	newUser := models.ReturnUser{ID: user.ID, Email: user.Email, Username: user.Username, Goal: user.Goal, UserRole: user.UserRole}
	if user.ImageUrl != "" {
		newUser.ImageUrl = os.Getenv("IMAGE_PATH") + user.ImageUrl
	}
	return newUser
}

func CreateToken(user models.User) (string, error) {
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["role"] = user.UserRole
	atClaims["id"] = user.ID
	atClaims["email"] = user.Email
	atClaims["exp"] = time.Now().Add(time.Hour * 72)

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}

	return token, nil
}

func isTokenValid(cookieToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(cookieToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil

	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return token, nil
}

func GetCurrentUserRole(cookieToken string) (models.Role, error) {
	token, err := isTokenValid(cookieToken)
	if err != nil {
		fmt.Println(err)
		return models.RoleNone, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return models.Role(claims["role"].(string)), nil
	}
	return models.RoleNone, nil
}

func GetUserByContext(c *gin.Context) (models.User, error) {
	tokenString, _ := c.Cookie("token")
	token, err := isTokenValid(tokenString)

	if err != nil {
		fmt.Println(err)
		return models.User{}, err
	}

	// get token Claims
	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		user, _, err := GetUserById(claims["id"].(string))
		if err != nil {
			return models.User{}, err
		}
		return user, nil
	}
	return models.User{}, nil

}

func SaveProfileImage(header *multipart.FileHeader, file multipart.File, profileUuid string) (string, error) {

	extension := filepath.Ext(header.Filename)
	filename := profileUuid + extension
	filePath := "profile/" + filename

	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}
	return filePath, nil
}
