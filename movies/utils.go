// Package movies provides a collection of video extensions
// and a utility function to check if given path is a video

package movies

import (
	"path"
	"path/filepath"
	"strings"
)

var extensions = []string{
	"3g2",
	"3gp",
	"aaf",
	"asf",
	"avchd",
	"avi",
	"flv",
	"m2v",
	"m4p",
	"m4v",
	"mkv",
	"mov",
	"mp2",
	"mp4",
	"mpg",
	"ogv",
	"vob",
	"webm",
	"wmv",
}

// Extensions is the extensions for different video types
var Extensions map[string]bool

func init() {
	Extensions = map[string]bool{}
	for _, ext := range extensions {
		Extensions[ext] = true
	}
}

// Is checks if a path is a video
func Is(p string) bool {
	ext := strings.TrimLeft(path.Ext(p), ".")
	return Extensions[strings.ToLower(ext)]
}

func FileName(p string) string {
	base := filepath.Base(p)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	return name
}
