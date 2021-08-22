package movies

import (
	"github.com/raspi/dirscanner"
)

var (
	AllMovies  = []*Movie{}
	AllSources = []*Source{}
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
	return db.AddMoviesToDB(AllSources)
}
