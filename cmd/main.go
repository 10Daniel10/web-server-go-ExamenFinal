package main

import (
	"fmt"
	"github.com/10Daniel10/web-server-go-ExamenFinal/internal/appointment"
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

	// Initialize and inject dependencies
	// Dentists
	dentistRepository := database.NewDentistRepository(db)
	dentistService := dentist.NewService(dentistRepository)
	dentistController := handler.NewDentistHandler(dentistService)

	// Patients
	patientRepository := database.NewPatientRepository(db)
	patientService := patient.NewService(patientRepository)
	patientController := handler.NewPatientHandler(patientService)

	// Appointments
	appointmentRepository := database.NewOtherAppointmentRepository(db)
	appointmentService := appointment.NewService(appointmentRepository)
	appointmentController := handler.NewAppointmentHandler(appointmentService, patientService, dentistService)

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
		// Configure routes
		dentistGroup.GET("", dentistController.GetAll)
		dentistGroup.GET("/q", dentistController.GetByLicense)
		dentistGroup.GET("/:id", dentistController.GetById)
		dentistGroup.POST("", authKeys.Validate, dentistController.Create)
		dentistGroup.PUT("/:id", authKeys.Validate, dentistController.Update)
		dentistGroup.PATCH("/:id", authKeys.Validate, dentistController.Patch)
		dentistGroup.DELETE("/:id", authKeys.Validate, dentistController.Delete)
	}

	patientGroup := baseGroup.Group("/patients")
	{
		// Configure routes
		patientGroup.GET("", patientController.GetAll)
		patientGroup.GET("/q", patientController.GetByDNI)
		patientGroup.GET("/:id", patientController.GetById)
		patientGroup.POST("", patientController.Create)
		patientGroup.PUT("/:id", patientController.Update)
		patientGroup.PATCH("/:id", patientController.Patch)
		patientGroup.DELETE("/:id", patientController.Delete)
	}

	appointmentGroup := baseGroup.Group("/appointments")
	{
		// Configure routes
		appointmentGroup.GET("", appointmentController.GetAll)
		appointmentGroup.GET("/:id", appointmentController.GetById)
		appointmentGroup.GET("/q", appointmentController.GetByDNI)
		appointmentGroup.POST("", appointmentController.Create)
		appointmentGroup.PUT("/:id", appointmentController.Update)
		appointmentGroup.PATCH("/:id", appointmentController.Patch)
		appointmentGroup.DELETE("/:id", appointmentController.Delete)

	}

	err = router.Run(envConfig.Private.Host)
	if err != nil {
		panic(fmt.Sprintf("Error running server: %v", err))
		return
	}
}
