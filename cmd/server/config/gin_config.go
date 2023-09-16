package config

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	return router
}
