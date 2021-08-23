package movies

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/raspi/dirscanner"
)

var (
	AllMovies  = []*Movie{}
	AllSources = []*Source{}
	tmdbClient = resty.New()
	tmbdMovies = []*TMDBMovie{}
)

type MovieInterface interface {
	AddMoviesToDB(movies []*Source) error
}

func Trans(files []dirscanner.FileInformation, db MovieInterface) error {

	for _, file := range files {
		name := FileName(file.Path)
		AllSources = append(AllSources, &Source{
			Name:     name,
			FileSize: float64(file.Size),
			FilePath: file.Path,
			Poster:   "",
		})
	}

	go processInfoFromTMDB()

	return db.AddMoviesToDB(AllSources)
}

func InitTMDBClient() {
	tmdbClient.SetHeader("Accept", "application/json")
	tmdbClient.SetHeader("Content-Type", "application/json")
}

type TMDBMovie struct {
	PosterPath   string `json:"poster_path"`
	BackdropPath string `json:"backdrop_path"`
	Id           int    `json:"id"`
	ReleaseDate  string `json:"release_date"`
	Title        string `json:"title"`
	Overview     string `json:"overview"`
}

type TMDBResponse struct {
	Page   int         `json:"page"`
	Result []TMDBMovie `json:"results"`
}

func queryMovie(query string, resource Source) {

	res, err := tmdbClient.
		NewRequest().
		SetQueryParam("query", query).
		SetQueryParam("api_key", "b03fc189ffc92b2d084ec96a81b2fd51").
		Get("https://api.themoviedb.org/3/search/movie")
	if err != nil {
		log.Println(err)
		return
	}
	data := res.Body()
	defer res.RawBody().Close()

	if res.IsError() {
		log.Println(errors.New("状态码错误:" + res.Status()))
		return
	}

	var response = TMDBResponse{}
	err = json.Unmarshal(data, &response)
	if err != nil {
		log.Println(err)
		return
	}

	if len(response.Result) == 0 {
		movie := Movie{
			Name:    resource.Name,
			Poster:  "",
			Sources: []Source{},
		}
		AllMovies = append(AllMovies, &movie)
		movie.Sources = append(movie.Sources, resource)
		return
	}
	tm := response.Result[0]
	getMovieForTMDBMovie(tm, resource)
}

func getMovieForTMDBMovie(tm TMDBMovie, source Source) Movie {

	for _, m := range AllMovies {
		if m.TMBDId == tm.Id {
			return *m
		}
	}

	movie := Movie{
		Name:   tm.Title,
		Poster: tm.PosterPath,
		Sources: []Source{
			source,
		},
	}
	defer func() {
		AllMovies = append(AllMovies, &movie)
	}()

	return movie
}

func processInfoFromTMDB() {

	for _, source := range AllSources {
		processForSource(source)
	}
}

func processForSource(resource *Source) {

	queryMovie(Parse(resource.Name), *resource)
}
