package main

import (
	"os"
	"path/filepath"

	"github.com/damoncoo/media-server/database"
	"github.com/prometheus/common/log"
)

const (
	appDirName = "media-server"
	dbFileName = "media-server.db"
)

func getAppDir() string {
	return filepath.Join(cacheDir, appDirName)
}

func getDBPath() string {
	return filepath.Join(getAppDir(), dbFileName)
}

func createDB(dbFile string) {

	// create working app dir
	appDirPath := getAppDir()
	if _, err := os.Stat(appDirPath); os.IsNotExist(err) {
		err = os.MkdirAll(appDirPath, 0777)
		if err != nil {
			panic("Unable to create App Dir on " + appDirPath)
		}
		log.Infoln("Created App dir at", appDirPath)
	}

	_, err := os.Stat(dbFile)
	if !os.IsNotExist(err) {
		log.Infoln("Obsolete DB detected, removing...")
		if err = os.Remove(dbFile); err != nil {
			panic("Unable removing obsolete DB")
		}
	}

	_, err = os.Create(dbFile)
	if err != nil {
		log.Errorln("Unable to init db file", err)
		os.Exit(1)
	}

	log.Infoln("DB initialized at", dbFile)

	database.Init(dbFile)

}
