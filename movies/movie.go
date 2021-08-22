package movies

type Source struct {
	Id       int     `xorm:"id not null pk autoincr INT(10)" json:"id" form:"id"`
	Name     string  `json:"name" xorm:"name" form:"name"`
	Poster   string  `json:"poster" xorm:"poster"`
	FileSize float64 `json:"fileSize" xorm:"file_size"`
	FilePath string  `json:"filePath" xorm:"file_path"`
}

type Movie struct {
	Id     int     `xorm:"id not null pk autoincr INT(10)" json:"id" form:"id"`
	Name   string  `json:"name" xorm:"name" form:"name"`
	Movies []Movie `json:"movies" xorm:"movies" form:"movies"`
}

func (Movie) TableName() string {

	return "movie"
}
