package service

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
	// find user with email
	err := usersCollection.FindOne(database.Ctx, bson.M{"email": email}).Decode(&user)

	return user, len(user.UUID) != 0, err
}
func GetUserByUsername(username string, usersCollection *mongo.Collection) (models.User, bool, error) {
	var user models.User
	// find user with email
	err := usersCollection.FindOne(database.Ctx, bson.M{"username": username}).Decode(&user)

	return user, len(user.UUID) != 0, err
}

func GetUserById(id string, usersCollection *mongo.Collection) (models.User, bool, error) {
	var user models.User
	// find user on UUID
	err := usersCollection.FindOne(database.Ctx, bson.M{"_id": id}).Decode(&user)
	return user, len(user.UUID) != 0, err
}

func UserToResultUser(user models.User) models.ReturnUser {
	newUser := models.ReturnUser{UUID: user.UUID, Email: user.Email, Username: user.Username, Goal: user.Goal}
	if user.ImageUrl != "" {
		newUser.ImageUrl = os.Getenv("IMAGE_PATH") + user.ImageUrl
	}
	return newUser
}

func CreateToken(user models.User) (string, error) {
	var err error
	// set new claims
	atClaims := jwt.MapClaims{}
	atClaims["role"] = user.UserRole
	atClaims["uuid"] = user.UUID
	atClaims["email"] = user.Email
	atClaims["exp"] = time.Now().Add(time.Hour * 72)

	// creating access token with claims
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}

	return token, nil
}

func isTokenValid(cookieToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(cookieToken, func(token *jwt.Token) (interface{}, error) {
		// make sure that the token method conform to "SigningMethodHMAC"
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

	// get token Claims
	claims, ok := token.Claims.(jwt.MapClaims)
	fmt.Println(claims)
	if ok && token.Valid {
		return models.Role(claims["role"].(string)), nil
	}
	return models.RoleNone, nil
}

func GetCurrentUserEmail(cookieToken string) (string, error) {
	token, err := isTokenValid(cookieToken)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// get token Claims
	claims, ok := token.Claims.(jwt.MapClaims)
	fmt.Println(claims)
	if ok && token.Valid {
		return claims["email"].(string), nil
	}
	return "", nil
}

func GetCurrentUserUuid(cookieToken string) (string, error) {
	token, err := isTokenValid(cookieToken)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// get token Claims
	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		return claims["uuid"].(string), nil
	}
	return "", nil
}

func GetUserByToken(c *gin.Context) (string, error) {
	tokenString, _ := c.Cookie("token")

	uuid, err := GetCurrentUserUuid(tokenString)
	fmt.Println("uuid" + uuid)
	return uuid, err

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
