package main

import (
	"fmt"
	"net/http"

	"github.com/10Daniel10/web-server-go-ExamenFinal/cmd/server/config"
	"github.com/10Daniel10/web-server-go-ExamenFinal/cmd/server/external/database"
	"github.com/10Daniel10/web-server-go-ExamenFinal/cmd/server/handler"
	"github.com/10Daniel10/web-server-go-ExamenFinal/cmd/server/middleware"
	"github.com/10Daniel10/web-server-go-ExamenFinal/internal/dentist"
	"github.com/10Daniel10/web-server-go-ExamenFinal/internal/patient"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

//	@title			Dental Clinic API
//	@version		2.0
//	@description	API of Dental Clinic
//	@termsOfService	http://swagger.io/terms/

// @tag.name	Patient
// @tag.description Patient operations for managing Patient
// @tag.docs.url http://swagger.io/terms/
// @tag.docs.description Patient operations for managing Patient

//	@accept		json
//	@produce	json

//	@schemes	http https

//	@contact.name	Gabriela Cecilia Calicanton
//	@contact.url	http://www.swagger.io/support
//	@contact.email	gabriela.calicanton@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@securityDefinitions.apikey	Bearer
//	@in							header
//	@name						Authorization
//	@description				Add Bearer token here, like this: Bearer {token}

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/

func main() {
	var err error

	err = godotenv.Load(".env")
	if err != nil {
		panic(fmt.Sprintf("Error loading .env file: %v", err))
		return
	}

	envConfig, err := config.NewEnvConfig("local")
	if err != nil {
		panic(fmt.Sprintf("Error loading config: %v", err))
		return
	}

	db, err := database.Connect(database.ConnectionParams{
		User:     envConfig.Private.DBUser,
		Password: envConfig.Private.DBPass,
		Host:     envConfig.Private.DBHost,
		Port:     envConfig.Private.DBPort,
		Database: envConfig.Private.DBName,
	})
	if err != nil {
		panic(fmt.Sprintf("Error connecting to database: %v", err))
		return
	}

	router := config.SetupRouter()
	{
		// Define global behavior
		router.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Not found"})
		})

		router.NoMethod(func(c *gin.Context) {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "Method not allowed"})
		})
	}
	baseGroup := router.Group(envConfig.Private.BasePath)
	{
		baseGroup.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
	}

	//Auth middleware
	authKeys := middleware.NewAuthKeys(envConfig.Private.SecretKey, envConfig.Public.PubKey)

	_ = baseGroup.Group("/docs")
	{
	}

	dentistGroup := baseGroup.Group("/dentists")
	{
		// Initialize and inject dependencies
		repository := database.NewDentistRepository(db)
		service := dentist.NewService(repository)
		controller := handler.NewDentistHandler(service)

		// Configure routes
		dentistGroup.GET("", controller.GetAll)
		dentistGroup.GET("/q", controller.GetByLicense)
		dentistGroup.GET("/:id", controller.GetById)
		dentistGroup.POST("", authKeys.Validate, controller.Create)
		dentistGroup.PUT("/:id", authKeys.Validate, controller.Update)
		dentistGroup.PATCH("/:id", authKeys.Validate, controller.Patch)
		dentistGroup.DELETE("/:id", authKeys.Validate, controller.Delete)
	}

	patientGroup := baseGroup.Group("/patients")
	{
		// Initialize and inject dependencies
		repository := database.NewPatientRepository(db)
		service := patient.NewService(repository)
		controller := handler.NewPatientHandler(service)

		// Configure routes
		patientGroup.GET("", controller.GetAll)
		patientGroup.GET("/q", controller.GetByDNI)
		patientGroup.GET("/:id", controller.GetById)
		patientGroup.POST("", controller.Create)
		patientGroup.PUT("/:id", controller.Update)
		patientGroup.PATCH("/:id", controller.Patch)
		patientGroup.DELETE("/:id", controller.Delete)
	}

	err = router.Run(envConfig.Private.Host)
	if err != nil {
		panic(fmt.Sprintf("Error running server: %v", err))
		return
	}
}
