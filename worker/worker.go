package worker

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fengdu/billconverter/converter"
	"github.com/fengdu/billconverter/input"
	"github.com/fengdu/billconverter/output"
)

// Start get files form src, then write csv to destination
func Start(src, destination string) {
	os.RemoveAll(destination)
	os.MkdirAll(destination, 0755)

	files, err := ioutil.ReadDir(src)
	if err != nil {
		log.Fatalln("ERROR: ReadDir: ", err)
	}

	var waitGroup sync.WaitGroup

	waitGroup.Add(len(files))

	for _, f := range files {
		go func(f os.FileInfo) {
			defer waitGroup.Done()
			if !f.IsDir() && strings.HasSuffix(f.Name(), ".txt") {
				_, err := process(f.Name(), src, destination)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Printf("INFO: %s convert successed.\n", f.Name())
			}

		}(f)
	}

	waitGroup.Wait()
	fmt.Println("INFO: all file convert successed.")
}

func process(filename, src, destination string) (string, error) {
	content, err := input.RetriveBillContent(src + "/" + filename)
	if err != nil {
		return "", fmt.Errorf("ERROR: read: %s", filename)
	}

	header, err := converter.GetBillBaseInfo(content)
	if err != nil {
		return "", fmt.Errorf("ERROR: GetBillBaseInfo: %s", filename)
	}

	now := time.Now()
	shortT := now.Format("20060102")
	longT := now.Format("20060102150405")

	accountNo := header["Account No"]

	tradeConfirmation := converter.GetTradeConfirmation(content)
	destFilename := fmt.Sprintf("%s_WANDA_SHTrades_%s_%s.csv", accountNo, shortT, longT)
	filepath := destination + "/" + destFilename
	if err := output.Write(filepath, tradeConfirmation); err != nil {
		return "", fmt.Errorf("ERROR: write: tradeConfirmation: %s", filename)
	}

	return filepath, nil
}
