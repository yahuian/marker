package utils

import (
	"regexp"
	"strings"
)

// BUG one line maybe have many images
var imageRegex = regexp.MustCompile(`!\[.*\]\((.*)\)`)

func ParseMarkdownImage(text string) string {
	list := imageRegex.FindAllStringSubmatch(text, 1)
	for _, v := range list {
		if len(v) != 2 {
			continue
		}
		image := v[1]
		// skip online image
		if strings.HasPrefix(image, "http://") || strings.HasPrefix(image, "https://") {
			continue
		}
		return image
	}
	return ""
}
