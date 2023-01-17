package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"io/ioutil"
	"lauf-du-sau/database"
	"lauf-du-sau/models"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func SaveImage(header *multipart.FileHeader, file multipart.File, folder string) (string, error) {
	now := time.Now()
	extension := filepath.Ext(header.Filename)
	filename := uuid.New().String() + "-" + fmt.Sprintf("%v", now.Unix()) + extension
	filePath := folder + filename

	out, err := os.Create(folder + filename)
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

func GetRun(id string) (models.Run, error) {
	runCollection := database.InitRunCollection()
	var run models.Run
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Run{}, err
	}
	err = runCollection.FindOne(database.Ctx, bson.M{"_id": objectId}).Decode(&run)
	if err != nil {
		return models.Run{}, err
	}

	return run, nil

}

func OcrImage(imageUrl string) {

	apiUrl := "https://microsoft-computer-vision3.p.rapidapi.com/ocr?detectOrientation=true&language=de"
	payload := strings.NewReader("{\n  \"url\": \"" + imageUrl + "\"\n}")

	req, _ := http.NewRequest("POST", apiUrl, payload)

	req.Header.Add("content-type", "application/json")
	req.Header.Add("X-RapidAPI-Key", "87849f7b4dmsh7ff24747209bbddp19f9c7jsnc1288aaf529a")
	req.Header.Add("X-RapidAPI-Host", "microsoft-computer-vision3.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println("------------------------------------")
	fmt.Println(res)
	fmt.Println(string(body))
	fmt.Println("------------------------------------")

}

func HasPermission(runId string, c *gin.Context) (bool, models.Run) {
	var run models.Run

	user, err := GetUserByContext(c)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return false, run
	}

	dbRun, err := GetRun(runId)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return false, run
	}
	if dbRun.UserID != user.ID && user.UserRole != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to delete this run"})
		return false, run
	}
	return true, dbRun
}

func FormatRun(run models.Run) models.RunResponse {
	var result models.RunResponse
	user, _, _ := GetUserById(run.UserID.Hex())
	result.User = UserToResultUser(user)
	result.Run = run
	result.Run.Url = os.Getenv("IMAGE_PATH") + run.Url

	return result

}
