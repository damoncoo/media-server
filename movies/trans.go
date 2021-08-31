package movies

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/go-resty/resty/v2"
)

var (
	AllMovies      = []*Movie{}
	AllSources     = []*Source{}
	tmdbClient     = resty.New()
	movieId    int = 1
)

type MovieInterface interface {
	AddSourcesToDB(sources []*Source) error
	AddMoviesToDB(sources []*Movie) error
}

func Trans(files []*Source, db MovieInterface) error {

	AllSources = files
	go processInfoFromTMDB(db)

	return nil
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
	var tm *TMDBMovie
	if len(response.Result) != 0 {
		tm = &(response.Result[0])
	}
	movie := getMovieForTMDBMovie(tm, resource)
	movie.Sources = append(movie.Sources, resource)
}

func getMovieForTMDBMovie(tm *TMDBMovie, source Source) *Movie {

	var movie Movie
	if tm != nil {
		for _, m := range AllMovies {
			if m.TMBDId == tm.Id {
				return m
			}
		}

		movie = Movie{
			Name:    tm.Title,
			Poster:  tm.PosterPath,
			Sources: []Source{},
			TMBDId:  tm.Id,
			DirPath: source.DirPath,
			Id:      movieId,
		}
	} else {
		movie = Movie{
			Name:    source.Name,
			Poster:  source.Poster,
			Sources: []Source{},
			DirPath: source.DirPath,
			Id:      movieId,
		}
		movieId += 1
	}

	defer func() {
		AllMovies = append(AllMovies, &movie)
	}()

	return &movie
}

func processInfoFromTMDB(db MovieInterface) {

	for _, source := range AllSources {
		processForSource(source)
	}
	_ = db.AddMoviesToDB(AllMovies)
}

func processForSource(resource *Source) {

	queryMovie(Parse(resource.Name), *resource)
}
