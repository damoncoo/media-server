package database

import (
	"time"

	"github.com/damoncoo/media-server/movies"
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

type MovieDAO struct {
	ORMEngine *xorm.Engine
}

var (
	DAO = MovieDAO{}
)

func (d MovieDAO) AddSourcesToDB(sources []*movies.Source) error {
	_, err := DAO.ORMEngine.Table("source").Insert(sources)
	return err
}

func (d MovieDAO) AddMoviesToDB(movies []*movies.Movie) error {
	_, err := DAO.ORMEngine.Table("movie").Insert(movies)
	return err
}

func Init(dbFile string) {

	Engine, err := xorm.NewEngine("sqlite3", dbFile)
	if err != nil {
		panic(err)
	}
	Engine.DB().SetConnMaxIdleTime(time.Second * 10)
	Engine.DB().SetMaxOpenConns(100)
	DAO.ORMEngine = Engine

	err = Engine.Sync2(new(movies.Movie), new(movies.Source))
	if err != nil {
		panic(err)
	}
}
