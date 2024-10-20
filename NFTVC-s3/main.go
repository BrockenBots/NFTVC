package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// На по***
const storagePath = "./storage"

type FileInfo struct {
	ID        string `bson:"_id,omitempty"`
	FileName  string `bson:"filename"`
	UploadURL string `bson:"upload_url"`
}

func connectMongo() (*mongo.Client, *mongo.Collection) {
	connect, err := mongo.Connect(context.Background(), options.Client().
		ApplyURI("mongodb://user:password@mongodb:27017/"))
	if err != nil {
		panic(err)
	}
	err = connect.Ping(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	client := connect.Database("S3").Collection("files")
	return connect, client
}

func generateFileName() string {
	uuid, _ := uuid.NewV7()
	return uuid.String()
}

func uploadFile(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.String(http.StatusBadRequest, "Не удалось получить файл")
	}
	src, err := file.Open()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Не удалось открыть файл")
	}
	defer src.Close()

	fileID := fmt.Sprintf("%s.jpg", generateFileName())
	filePath := filepath.Join(storagePath, fileID)

	dst, err := os.Create(filePath)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Не удалось создать файл на сервере")
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return c.String(http.StatusInternalServerError, "Ошибка при копировании файла")
	}

	downloadURL := fmt.Sprintf("http://localhost:8083/files/%s", fileID)

	client, collection := connectMongo()
	defer client.Disconnect(context.TODO())

	fileInfo := FileInfo{
		ID:        fileID,
		FileName:  file.Filename,
		UploadURL: downloadURL,
	}

	_, err = collection.InsertOne(context.TODO(), fileInfo)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Не удалось сохранить информацию о файле в базу данных")
	}

	return c.String(http.StatusOK, downloadURL)
}

func getFileLink(c echo.Context) error {
	fileID := c.Param("id")

	client, collection := connectMongo()
	defer client.Disconnect(context.TODO())

	var fileInfo FileInfo
	err := collection.FindOne(context.TODO(), bson.M{"_id": fileID}).Decode(&fileInfo)
	if err != nil {
		return c.String(http.StatusNotFound, "Файл не найден")
	}

	return c.String(http.StatusOK, fileInfo.UploadURL)
}

func downloadFile(c echo.Context) error {
	fileID := c.Param("id")

	filePath := filepath.Join(storagePath, fileID)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.String(http.StatusNotFound, "Файл не найден")
	}

	return c.File(filePath)
}

func deleteFile(c echo.Context) error {
	fileID := c.Param("id")

	client, collection := connectMongo()
	defer client.Disconnect(context.TODO())

	_, err := collection.DeleteOne(context.TODO(), bson.M{"_id": fileID})
	if err != nil {
		return c.String(http.StatusInternalServerError, "Ошибка при удалении информации о файле из базы данных")
	}

	filePath := filepath.Join(storagePath, fileID)

	err = os.Remove(filePath)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Ошибка при удалении файла")
	}

	return c.String(http.StatusOK, fmt.Sprintf("Файл удален: %s", fileID))
}

func main() {
	err := os.MkdirAll(storagePath, os.ModePerm)
	if err != nil {
		log.Fatalf("Ошибка создания директории для хранения файлов: %v", err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/upload", uploadFile)
	e.GET("/files/:id", downloadFile)
	e.GET("/link/:id", getFileLink)
	e.DELETE("/files/:id", deleteFile)

	log.Println("Сервер запущен на порту :8083")
	err = e.Start(":8083")
	if err != nil {
		log.Fatal(err)
	}
}
