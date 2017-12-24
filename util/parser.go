package util

import (
	"fmt"
	"regexp"
)

// ParseSegment parse segment by title
func ParseSegment(content, title string) string {
	var matched string
	pattern := fmt.Sprintf(`(?P<first>(\b%s.*\t*[\r\n]))(.|\n)*?(-+\t*[\r\n])`, title)
	r := regexp.MustCompile(pattern)

	m := r.FindStringSubmatch(content)
	if len(m) > 0 {
		matched = m[0]
	}

	return matched
}
