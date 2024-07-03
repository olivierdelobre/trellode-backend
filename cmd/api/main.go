package main

import (
	docs "trellode-go/docs"
	"trellode-go/internal/api"
	"trellode-go/internal/middlewares"

	"trellode-go/internal/utils/config"

	"github.com/gin-gonic/gin"
)

// @title           			Lists API
// @version         			1.0
// @description     			This is the Lists API
// @contact.name    			Contact ISCS-IAM
// @contact.email   			idev-md@groupes.epfl.ch
// @host            			api.epfl.ch
// @Security BasicAuth
// @securityDefinitions.basic	BasicAuth
func main() {
	docs.SwaggerInfo.Title = "Lists API"

	c := config.GetConfig()
	// Get logger from config
	log := c.Log

	// Get db from config
	db := c.Db

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middlewares.CorsMiddleware(log))
	r.Use(middlewares.AuthenticationMiddleware(db, log))
	r.Use(middlewares.LoggingMiddleware(log))

	s := api.NewServer(db, r, log)

	s.Routes()

	err := r.Run()
	if err != nil {
		return
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
