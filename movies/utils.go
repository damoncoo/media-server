// Package movies provides a collection of video extensions
// and a utility function to check if given path is a video

package movies

import (
	"path"
	"path/filepath"
	"reflect"
	"regexp"
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

var patterns = []struct {
	name string
	// Use the last matching pattern. E.g. Year.
	last bool
	kind reflect.Kind
	// REs need to have 2 sub expressions (groups), the first one is "raw", and
	// the second one for the "clean" value.
	// E.g. Epiode matching on "S01E18" will result in: raw = "E18", clean = "18".
	re *regexp.Regexp
}{
	{"season", false, reflect.Int, regexp.MustCompile(`(?i)(s?([0-9]{1,2}))[ex]`)},
	{"episode", false, reflect.Int, regexp.MustCompile(`(?i)([ex]([0-9]{2})(?:[^0-9]|$))`)},
	{"episode", false, reflect.Int, regexp.MustCompile(`(-\s+([0-9]{1,})(?:[^0-9]|$))`)},
	{"year", true, reflect.Int, regexp.MustCompile(`\b(((?:19[0-9]|20[0-9])[0-9]))\b`)},

	{"resolution", false, reflect.String, regexp.MustCompile(`\b(([0-9]{3,4}p))\b`)},
	{"quality", false, reflect.String, regexp.MustCompile(`(?i)\b(((?:PPV\.)?[HP]DTV|(?:HD)?CAM|B[DR]Rip|(?:HD-?)?TS|(?:PPV )?WEB-?DL(?: DVDRip)?|HDRip|DVDRip|DVDRIP|CamRip|W[EB]BRip|BluRay|DvDScr|telesync))\b`)},
	{"codec", false, reflect.String, regexp.MustCompile(`(?i)\b((xvid|[hx]\.?26[45]))\b`)},
	{"audio", false, reflect.String, regexp.MustCompile(`(?i)\b((MP3|DD5\.?1|Dual[\- ]Audio|LiNE|DTS|AAC[.-]LC|AAC(?:\.?2\.0)?|AC3(?:\.5\.1)?))\b`)},
	{"region", false, reflect.String, regexp.MustCompile(`(?i)\b(R([0-9]))\b`)},
	{"size", false, reflect.String, regexp.MustCompile(`(?i)\b((\d+(?:\.\d+)?(?:GB|MB)))\b`)},
	{"website", false, reflect.String, regexp.MustCompile(`^(\[ ?([^\]]+?) ?\])`)},
	{"language", false, reflect.String, regexp.MustCompile(`(?i)\b((rus\.eng|ita\.eng))\b`)},
	{"sbs", false, reflect.String, regexp.MustCompile(`(?i)\b(((?:Half-)?SBS))\b`)},
	{"container", false, reflect.String, regexp.MustCompile(`(?i)\b((MKV|AVI|MP4))\b`)},

	{"group", false, reflect.String, regexp.MustCompile(`\b(- ?([^-]+(?:-={[^-]+-?$)?))$`)},

	{"extended", false, reflect.Bool, regexp.MustCompile(`(?i)\b(EXTENDED(:?.CUT)?)\b`)},
	{"hardcoded", false, reflect.Bool, regexp.MustCompile(`(?i)\b((HC))\b`)},
	{"proper", false, reflect.Bool, regexp.MustCompile(`(?i)\b((PROPER))\b`)},
	{"repack", false, reflect.Bool, regexp.MustCompile(`(?i)\b((REPACK))\b`)},
	{"widescreen", false, reflect.Bool, regexp.MustCompile(`(?i)\b((WS))\b`)},
	{"unrated", false, reflect.Bool, regexp.MustCompile(`(?i)\b((UNRATED))\b`)},
	{"threeD", false, reflect.Bool, regexp.MustCompile(`(?i)\b((3D))\b`)},
}

func Parse(filename string) string {

	var startIndex, endIndex = 0, len(filename)
	cleanName := strings.Replace(filename, "_", " ", -1)
	for _, pattern := range patterns {
		matches := pattern.re.FindAllStringSubmatch(cleanName, -1)
		if len(matches) == 0 {
			continue
		}
		matchIdx := 0
		if pattern.last {
			matchIdx = len(matches) - 1
		}

		index := strings.Index(cleanName, matches[matchIdx][1])
		if index == 0 {
			startIndex = len(matches[matchIdx][1])
			//fmt.Printf("    startIndex moved to %d [%q]\n", startIndex, filename[startIndex:endIndex])
		} else if index < endIndex {
			endIndex = index
		}
	}

	// Start process for title
	//fmt.Println("  title: <internal>")
	raw := strings.Split(filename[startIndex:endIndex], "(")[0]
	cleanName = raw
	if strings.HasPrefix(cleanName, "- ") {
		cleanName = raw[2:]
	}
	if strings.ContainsRune(cleanName, '.') && !strings.ContainsRune(cleanName, ' ') {
		cleanName = strings.Replace(cleanName, ".", " ", -1)
	}
	cleanName = strings.Replace(cleanName, "_", " ", -1)
	cleanName = strings.TrimSpace(cleanName)
	return cleanName
}
