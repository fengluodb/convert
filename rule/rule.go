package rule

import (
	"regexp"
	"strings"
)

func TxtTitleRegexp() *regexp.Regexp {
	titleRegexp, _ := regexp.Compile("第[一二三四五六七八九零〇一二三四五六七八九十百千万a-zA-Z0-9]{1,7}[章节卷集部篇回].*")

	return titleRegexp
}

func IsTxt(filename string) bool {
	tmp := strings.Split(filename, ".")
	if len(tmp) == 2 && tmp[1] == "txt" {
		return true
	}
	return false
}
