package main

import (
	"SlackPlatform/controller"
	"SlackPlatform/middleware"
	"SlackPlatform/models"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

//DATABASE_URL

func main() {
	db := connectToGormDB()
	defer db.Close()

	models.SetDatabase(db)
	controller.StartupControllers(db)

	//Middleware
	middlewareHandlers := new(middleware.GzipMiddleware)
	middlewareHandlers.Next = new(middleware.TeamScope)

	err := http.ListenAndServe(determineListenAddress(), middlewareHandlers)
	if err != nil {
		fmt.Println(fmt.Errorf("Error listening to 8080, %v", err))
		panic(err)
	}
}

func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return ":8080", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
}

func determineDBConnection() string {
	dbConnection := os.Getenv("DATABASE_URL")
	const connStr = "user=goUser dbname=barista password=qwe123 sslmode=disable"
	if dbConnection == "" {
		return connStr
	}
	return dbConnection
}

func connectToGormDB() *gorm.DB {
	connStr := determineDBConnection()
	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
