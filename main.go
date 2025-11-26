package main

import (
	"crud_alumni/config"
	"crud_alumni/database"
	"crud_alumni/route"
	"log"

	_ "crud_alumni/docs"

	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title CRUD Alumni API
// @version 1.0
// @description API untuk mengelola data alumni, pekerjaan, dan upload file (foto & sertifikat)
// @host localhost:3000
// @BasePath /api
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	config.LoadEnv()
	config.InitLogger()
	database.ConnectDB()

	app := config.App()

	// route setup
	route.SetupRoutes(app)

	// Swagger route
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Static uploads
	app.Static("/uploads", "./uploads")

	port := config.GetEnv("APP_PORT", "3000")
	log.Fatal(app.Listen(":" + port))
}
