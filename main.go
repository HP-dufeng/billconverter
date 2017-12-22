package main

import (
	"time"
	"bufio"
	"os"
	"strings"
	"fmt"
	"regexp"
	"io/ioutil"
	"encoding/csv"
	"log"
)

var src string = "./src"
var destination string = "./dst"

func main() {
	files, err := ioutil.ReadDir(src)
	if err != nil {
		log.Fatalln("ERROR: ReadDir: ", err)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		filename := f.Name()
		content, err := getFileContent(filename)
		if err != nil {
			log.Fatalln(fmt.Sprintf("ERROR: getFileContent: %s", filename), err)
		}

		segment := parseSegment(content, "Financial Situation")
		if len(segment) <=0 {
			continue
		}
		destFileName := fmt.Sprintf("WANDA_SHBalances_%s_%s.csv", time.Now().Format("20060102"), time.Now().Format("20060102150405"))
		filepath := destination + "/" + destFileName
		toCsv(filepath, segment)

	}

	fmt.Println("INFO: all file convert successed.")
}

	
// func check(e error) {
//     if e != nil {
// 		fmt.Println(e)
// 		os.Exit(1)
//     }
// }

// func from(dir string) ([]os.FileInfo, error) {
// 	files, err := ioutil.ReadDir(dir)
// 	return files, err
// }

// func to(fileInfos []os.FileInfo, dir string) {
// 	for _, file := range fileInfos {

// 	}
// }

func getFileContent(filename string) (string, error) {
	var content string
	f, err := ioutil.ReadFile(src + "/" + filename)
	if err == nil {
		content = string(f)
		if content = strings.TrimSpace(content); strings.HasSuffix(content,"------") {
			content += "\r\n"
		} else {
			content += "\r\n\t-------\r\n"
		}
	}

	return content, err
}

func parseSegment(content, title string) string {
	var matched string
	pattern := fmt.Sprintf(`(?P<first>(\b%s.*\t*[\r\n]))(.|\n)*?(-+\t*[\r\n])`, title)
	r := regexp.MustCompile(pattern)

	m := r.FindStringSubmatch(content)
	if len(m) > 0 {
		matched = m[0]
	} 

	return matched
}

func toCsv(filepath, segment string) {
	// records := [][]string{
	// 	{"first_name", "last_name", "username"},
	// 	{"Rob", "Pike", "rob"},
	// 	{"Ken", "Thompson", "ken"},
	// 	{"Robert", "Griesemer", "gri"},
	// }
	f, err := os.Create(filepath)
	defer f.Close()
	if err !=nil {
		log.Fatalln("ERROR: toCsv: ", err)
	}

	w := csv.NewWriter(f)
	scanner := bufio.NewScanner(strings.NewReader(segment))
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), "|")
		if!strings.Contains(line, "|") {
			continue
		}
		
		if err := w.Write(strings.Split(line, "|")); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	w.Flush()
	if err := scanner.Err(); err != nil {
		log.Fatalln("reading standard input:", err)
	}
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}