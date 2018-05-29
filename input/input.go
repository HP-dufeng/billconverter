package input

import (
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// RetriveBillContent read bill content from file path
func RetriveBillContent(filepath string) (string, error) {
	var content string
	f, err := os.Open(filepath)
	defer f.Close()
	if err == nil {
		ts := simplifiedchinese.GBK.NewDecoder()
		r := transform.NewReader(f, ts)
		b, err := ioutil.ReadAll(r)
		if err != nil {
			return "", err
		}

		content = string(b)
		if content = strings.TrimSpace(content); strings.HasSuffix(content, "------") {
			content += "\r\n"
		} else {
			content += "\r\n\t-------\r\n"
		}
	}

	return content, err
}
