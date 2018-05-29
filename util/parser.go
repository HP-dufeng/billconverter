package util

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
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

// ParseMonthAndYear parse month & year from date format string, eq: "1812" or "812"
func ParseMonthAndYear(d string) (month string, year string, err error) {
	if !(len(d) == 3 || len(d) == 4) {
		err = errors.New("Paramer error, length must be 3 or 4")
		return
	}

	if len(d) == 3 {
		d = "1" + d
	}

	var t time.Time
	t, err = time.Parse("0601", d)
	if err != nil {
		return
	}

	month = strconv.Itoa(int(t.Month()))
	year = strconv.Itoa(t.Year())

	return
}
