package http

import (
	"crypto/md5"
	"fmt"
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
		movie.GET("/play/", Play)
	}

	log.Printf("Server started at 0.0.0.0:%s", port)

	for _, dir := range paths {
		app.Static("/static/"+Md5(dir), dir)
	}

	err := app.Run(":" + port)
	if err != nil {
		panic(err)
	}

}

// Md5 str 取MD5
func Md5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str
}
