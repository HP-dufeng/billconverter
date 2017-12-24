package converter

import (
	"bufio"
	"errors"
	"regexp"
	"strings"

	"github.com/fengdu/billconverter/util"
)

// GetBillBaseInfo parse bill base info segment
func GetBillBaseInfo(content string) (map[string]string, error) {
	s := util.ParseSegment(content, "Account No")
	s = strings.Replace(s, "：", ":", -1)

	result := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		line := strings.Trim(strings.TrimSpace(scanner.Text()), "|")
		if regexp.MustCompile(`^-+`).MatchString(line) {
			continue
		}
		p := regexp.MustCompile(`\s{2,}`).Split(line, -1)
		for _, f := range p {
			kv := strings.Split(f, ":")
			if len(kv) != 2 {
				return nil, errors.New("Parse bill base info errors")
			}
			result[kv[0]] = kv[1]
		}
	}

	return result, nil
}

// GetTradeConfirmation parse Trade Confirmation segment
func GetTradeConfirmation(content string) [][]string {
	s := util.ParseSegment(content, "Trade Confirmation")

	result := [][]string{}
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		line := strings.Trim(strings.TrimSpace(scanner.Text()), "|")
		if !strings.Contains(line, "|") {
			continue
		}

		result = append(result, strings.Split(line, "|"))
	}

	return result
}
