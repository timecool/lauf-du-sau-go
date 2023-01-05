package service

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
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
