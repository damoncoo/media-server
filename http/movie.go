package http

import (
	"strconv"
	"strings"

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

func Play(ctx *gin.Context) {

	type R struct {
		Key string `form:"key"`
	}

	var r R
	ctx.ShouldBindQuery(&r)

	key := r.Key
	m, s := ParseMovieKey(key)

	movie := movies.AllMovies[m]
	source := movie.Sources[s]
	ctx.File(source.FilePath)
}

func ParseMovieKey(key string) (int, int) {
	keys := strings.Split(key, "-")

	var m, s int
	if len(keys) == 2 {
		m, _ = strconv.Atoi(keys[0])
		m, _ = strconv.Atoi(keys[1])
	}
	return m, s
}
