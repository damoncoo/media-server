package http

import (
	"log"

	"github.com/gin-gonic/gin"
)

func Serve(port string, paths ...string) {

	app := gin.New()
	app.Use(gin.Recovery())
	app.Use(gin.Logger())
	app.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	movie := app.Group("/movie")
	{
		movie.GET("/", Movies)
	}

	log.Printf("Server started at 0.0.0.0:%s", port)

	for _, dir := range paths {
		app.Static("/static", dir)
	}

	err := app.Run(":" + port)
	if err != nil {
		panic(err)
	}

}
