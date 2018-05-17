package main

import (
	"SlackPlatform/controller"
	"SlackPlatform/middleware"
	"SlackPlatform/models"
	"fmt"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

const (
	connStr = "user=goUser dbname=barista password=qwe123 sslmode=disable"
)

func main() {
	db := connectToGormDB()
	defer db.Close()

	models.SetDatabase(db)
	controller.StartupControllers(db)

	//Middleware
	middlewareHandlers := new(middleware.GzipMiddleware)
	middlewareHandlers.Next = new(middleware.TeamScope)

	err := http.ListenAndServe(":8080", middlewareHandlers)
	if err != nil {
		fmt.Println(fmt.Errorf("Error listening to 8080, %v", err))
		panic(err)
	}
}

func connectToGormDB() *gorm.DB {
	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
