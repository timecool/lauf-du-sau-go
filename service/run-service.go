package service

import (
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
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

func GetUserFromRunUuid(uuid string) (string, error) {
	userCollection := database.InitUserCollection()

	o1 := bson.M{
		"$unwind": "$runs",
	}
	o2 := bson.M{
		"$match": bson.M{"runs._id": uuid},
	}
	o4 := bson.M{
		"$group": bson.M{
			"_id": "$_id",
		},
	}

	cursor, err := userCollection.Aggregate(database.Ctx, []bson.M{o1, o2, o4})
	var results []models.User
	if err = cursor.All(database.Ctx, &results); err != nil {
		return "", err
	}
	if err := cursor.Close(database.Ctx); err != nil {
		return "", err
	}

	if len(results) == 1 {
		return results[0].UUID, nil
	}
	return "", nil

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
