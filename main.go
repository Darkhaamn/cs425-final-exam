package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"prince-university/student-records/controller"
	"prince-university/student-records/database"
	"prince-university/student-records/repository"
	"prince-university/student-records/service"
)

func main() {
	_ = godotenv.Load()

	db := database.Connect()

	studentRepo := repository.NewStudentRepository(db)
	studentSvc := service.NewStudentService(studentRepo)
	studentCtl := controller.NewStudentController(studentSvc)

	r := gin.Default()
	api := r.Group("/api/v1")
	studentCtl.RegisterRoutes(api)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Prince University Student Records API listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
