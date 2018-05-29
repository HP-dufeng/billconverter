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
	// Clear dst folder
	if stat, err := os.Stat(destination); err == nil && stat.IsDir() {
		temp := fmt.Sprintf("_%v", time.Now().Unix())
		os.Rename(destination, temp)
		os.RemoveAll(temp)
	}
	if err := os.MkdirAll(destination, 0777); err != nil {
		fmt.Println("ERROR: MkdirAll: ", err)
		return
	}

	// Read bills from src folder
	files, err := ioutil.ReadDir(src)
	if err != nil {
		log.Fatalln("ERROR: ReadDir: ", err)
	}

	var waitGroup sync.WaitGroup

	waitGroup.Add(len(files))

	for _, f := range files {
		// Convert to csv file individually
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

func process(filename, src, destination string) ([]string, error) {
	content, err := input.RetriveBillContent(src + "/" + filename)
	if err != nil {
		return nil, fmt.Errorf("ERROR: read: %s", filename)
	}

	// Get bill base info
	bill, err := converter.GetBillBaseInfo(content)
	if err != nil {
		return nil, fmt.Errorf("ERROR: GetBillBaseInfo: %s: %v", filename, err)
	}

	now := time.Now()
	shortT := now.Format("20060102")
	longT := now.Format("20060102150405")

	// Convert segments to csv
	filepaths := []string{}

	var fp string
	fp, err = writeBalances(&content, destination, shortT, longT, bill)
	if err != nil {
		return nil, err
	}
	filepaths = append(filepaths, fp)

	fp, err = writePos(&content, destination, shortT, longT, bill)
	if err != nil {
		return nil, err
	}
	filepaths = append(filepaths, fp)

	fp, err = writeTrades(&content, destination, shortT, longT, bill)
	if err != nil {
		return nil, err
	}
	filepaths = append(filepaths, fp)

	return filepaths, nil
}

func writeBalances(content *string, destination, shortT, longT string, bill converter.BillBaseInfo) (string, error) {
	data := converter.GetBalances(bill, *content)
	filename := fmt.Sprintf("%s_WANDA_SHBalances_%s_%s.csv", bill.AccountNo, shortT, longT)
	filepath := destination + "/" + filename
	if err := output.Write(filepath, data); err != nil {
		return "", fmt.Errorf("ERROR: write: Balances: %s：%v", filename, err)
	}

	return filepath, nil
}

func writePos(content *string, destination, shortT, longT string, bill converter.BillBaseInfo) (string, error) {
	data := converter.GetPos(bill, *content)
	filename := fmt.Sprintf("%s_WANDA_SHPos_%s_%s.csv", bill.AccountNo, shortT, longT)
	filepath := destination + "/" + filename
	if err := output.Write(filepath, data); err != nil {
		return "", fmt.Errorf("ERROR: write: Pos: %s：%v", filename, err)
	}

	return filepath, nil
}

func writeTrades(content *string, destination, shortT, longT string, bill converter.BillBaseInfo) (string, error) {
	data := converter.GetTrades(bill, *content)
	filename := fmt.Sprintf("%s_WANDA_SHTrades_%s_%s.csv", bill.AccountNo, shortT, longT)
	filepath := destination + "/" + filename
	if err := output.Write(filepath, data); err != nil {
		return "", fmt.Errorf("ERROR: write: Trades: %s：%v", filename, err)
	}

	return filepath, nil
}
