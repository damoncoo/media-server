package http

import (
	"github.com/damoncoo/media-server/movies"
	"github.com/gin-gonic/gin"
)

func Movies(ctx *gin.Context) {

	ctx.JSON(200, movies.Response{
		Code:    200,
		Message: "ok",
		Data:    movies.AllMovies,
	})
}
