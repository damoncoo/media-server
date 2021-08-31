package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/damoncoo/media-server/conf"
	"github.com/damoncoo/media-server/database"
	"github.com/damoncoo/media-server/http"
	"github.com/damoncoo/media-server/movies"
	"github.com/damoncoo/media-server/scanner"
	"github.com/jaffee/commandeer"
	"gopkg.in/yaml.v2"
)

// Command Command
type Command struct {
	Port     int    `help:"Port Of App" flag:"p"`
	Config   string `help:"Path Of Config" flag:"c"`
	CreateDB bool   `help:"Wether create db file" flag:"C"`
}

func defaultCommand() *Command {
	return &Command{
		Port:     3000,
		Config:   "conf.yml",
		CreateDB: false,
	}
}

// Run Run
func (m *Command) Run() error {

	return parseConfig(m.Config)
}

var (
	command     *Command
	cacheDir, _ = os.UserCacheDir()
	AllSources  = []*movies.Source{}
	paths       = []string{}
)

// 解析参数
func parseArgv() {

	command = defaultCommand()
	err := commandeer.Run(command)
	if err != nil {
		panic(err)
	}

}

// 解析配置文件
func parseConfig(pathConfig string) error {

	data, err := ioutil.ReadFile(pathConfig)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &conf.Conf)
}

func main() {

	// parse argv
	parseArgv()

	//  init db
	db := getDBPath()
	createDB(db)

	// map all pathes to get all movie files
	for _, path := range conf.Conf.Path {
		paths = append(paths, path)
		files, err := scanner.FindAllMovies(path)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, file := range files {
			name := movies.FileName(file.Path)
			s := &movies.Source{
				Name:     name,
				FileSize: float64(file.Size),
				FilePath: file.Path,
				Poster:   "",
				DirPath:  path,
			}
			AllSources = append(AllSources, s)
		}
	}

	// sort all movies
	err := movies.Trans(AllSources, database.DAO)
	if err != nil {
		log.Println(err)
	}

	http.Serve(strconv.Itoa(command.Port), paths...)
}
