package scanner

import (
	"log"
	"runtime"

	"github.com/damoncoo/media-server/movies"
	"github.com/raspi/dirscanner"
)

var (
	myScanner   dirscanner.DirectoryScanner
	workerCount int
)

func init() {
	workerCount = runtime.NumCPU()
	myScanner = dirscanner.New()

	err := myScanner.Init(workerCount, validateFile)
	if err != nil {
		panic(err)
	}
}

//  Custom file validator
func validateFile(info dirscanner.FileInformation) bool {
	return info.Mode.IsRegular() && movies.Is(info.Path)
}

func FindAllMovies(path string) ([]dirscanner.FileInformation, error) {

	var files []dirscanner.FileInformation

	err := myScanner.Init(workerCount, validateFile)
	if err != nil {
		panic(err)
	}

	err = myScanner.ScanDirectory(path)
	if err != nil {
		panic(err)
	}
	defer myScanner.Close()

scanloop:
	for {
		select {

		case <-myScanner.Finished: // Finished getting file list
			log.Printf("Find all movies for path: %s", path)
			break scanloop

		case e, ok := <-myScanner.Errors: // Error happened, handle, discard or abort
			if ok {
				log.Printf(`got error: %v`, e)
				//s.Aborted <- true // Abort
			}

		case res, ok := <-myScanner.Results:
			if ok {
				files = append(files, res)
			}

		case <-myScanner.Information:
		}
	}

	return files, err
}

func FindAllMoviePathes(path string) ([]string, error) {

	files, err := FindAllMovies(path)
	paths := []string{}

	for _, info := range files {
		paths = append(paths, info.Path)
	}
	return paths, err
}
