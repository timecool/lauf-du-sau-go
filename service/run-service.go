package service

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
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
