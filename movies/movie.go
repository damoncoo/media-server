package movies

type Source struct {
	Id       int     `xorm:"id not null pk autoincr INT(10)" json:"id" form:"id"`
	Name     string  `json:"name" xorm:"name" form:"name"`
	Poster   string  `json:"poster" xorm:"poster"`
	FileSize float64 `json:"fileSize" xorm:"file_size"`
	FilePath string  `json:"filePath" xorm:"file_path"`
}

type Movie struct {
	Id      int      `json:"id" xorm:"id not null pk autoincr INT(10)" form:"id"`
	Name    string   `json:"name" xorm:"name" form:"name"`
	Poster  string   `json:"poster" xorm:"poster"`
	Sources []Source `json:"sources" xorm:"sources" form:"sources"`
	TMBDId  int      `json:"tmbd_id" xorm:"-" `
}

// Subtitle is Gorm model of subtitle
type Subtitle struct {
	Rid      int    `json:"rid" xorm:"rid"`
	DirPath  string `json:"dir_path" xorm:"dir_path"`
	BaseName string `json:"base_name" xorm:"base_name"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (Movie) TableName() string {

	return "movie"
}
