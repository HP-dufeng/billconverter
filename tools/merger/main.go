package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fengdu/billconverter/output"
)

func main() {
	mainflag := flag.String("main", "./dst", "主账单目录")
	subflag := flag.String("sub", "./dst_sub", "子帐单目录")
	destinationflag := flag.String("dst_merge", "./dst_merge", "合并文件目录")

	flag.Parse()

	if msg, err := checkDir(*mainflag, *subflag, *destinationflag); err != nil {
		fmt.Println(msg)
		return
	}

	mainBillPath := *mainflag
	subBillPath := *subflag
	destination := *destinationflag

	fmt.Println("主账单目录: ", mainBillPath)
	fmt.Println("子账单目录: ", subBillPath)
	fmt.Println("合并文件目录: ", destination)

	mergeBalances(mainBillPath, subBillPath, destination)
	mergePoses(mainBillPath, subBillPath, destination)
	mergeTrades(mainBillPath, subBillPath, destination)
}

func checkDir(mainBillPath, subBillPath, destination string) (string, error) {
	if _, err := os.Stat(mainBillPath); os.IsExist(err) {
		return "主账单目录不存在.", err
	}
	if _, err := os.Stat(subBillPath); os.IsExist(err) {
		return "子帐单目录不存在.", err
	}

	if _, err := os.Stat(destination); os.IsNotExist(err) {
		if err = os.MkdirAll(destination, 0777); err != nil {
			return "ERROR: 创建合并目录错误 ", err
		}
	}

	return "", nil
}

func getBillsInfo(pattern, mainBillPath, subBillPath string) (mainBillFileNames, subBillFileNames []string) {
	// Read bills from folder
	var mainBillFiles, subBillFiles []os.FileInfo
	var err error
	if len(mainBillPath) > 0 {
		if mainBillFiles, err = ioutil.ReadDir(mainBillPath); err != nil {
			fmt.Println("读取主账单目录错误", err)
		}

		for _, f := range mainBillFiles {
			if f.IsDir() {
				continue
			}
			if strings.Contains(f.Name(), pattern) {
				mainBillFileNames = append(mainBillFileNames, f.Name())
			}
		}
	}

	if len(subBillPath) > 0 {
		if subBillFiles, err = ioutil.ReadDir(subBillPath); err != nil {
			fmt.Println("读取子帐单目录错误", err)
		}

		for _, f := range subBillFiles {
			if f.IsDir() {
				continue
			}
			if strings.Contains(f.Name(), pattern) {
				subBillFileNames = append(subBillFileNames, f.Name())
			}
		}
	}

	return mainBillFileNames, subBillFileNames
}

func readCSV(filepath string) ([][]string, error) {
	f, err := os.Open(filepath)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	// Read File into a Variable
	lines, err := csv.NewReader(f).ReadAll()
	return lines, err

	// // Loop through lines & turn into object
	// for _, line := range lines {
	//     data := CsvLine{
	//         Column1: line[0],
	//         Column2: line[1],
	//     }
	//     fmt.Println(data.Column1 + " " + data.Column2)
	// }
}

func mergeBalances(mainBillPath, subBillPath, destination string) {
	mainBillFiles, subBillFiles := getBillsInfo("Balance", mainBillPath, subBillPath)

	if len(mainBillFiles) <= 0 {
		fmt.Println("Balances 文件为空")
		return
	}

	rows, err := readCSV(mainBillPath + "/" + mainBillFiles[0])
	if err != nil {
		fmt.Println("合并 Balances 错误")
		return
	}

	for _, filename := range subBillFiles {
		lines, err := readCSV(subBillPath + "/" + filename)
		if err != nil {
			fmt.Println("合并 Balances 错误")
			fmt.Println(err)
			break
		}

		for _, line := range lines[1:] {
			rows = append(rows, line)
		}
	}

	output.Write(destination+"/"+mainBillFiles[0], rows)

}

func mergePoses(mainBillPath, subBillPath, destination string) {
	mainBillFiles, subBillFiles := getBillsInfo("Pos", mainBillPath, subBillPath)

	if len(mainBillFiles) <= 0 {
		fmt.Println("Pos 文件为空")
		return
	}

	if len(subBillFiles) <= 0 {
		fmt.Println("Pos 文件为空")
		return
	}

	temp, _ := readCSV(subBillPath + "/" + subBillFiles[0])
	header := temp[0]

	rows := [][]string{header}

	for _, filename := range subBillFiles {
		lines, err := readCSV(subBillPath + "/" + filename)
		if err != nil {
			fmt.Println("合并 Pos 错误")
			fmt.Println(err)
			break
		}

		for _, line := range lines[1:] {
			tradedate := line[1]
			if t, err := time.Parse("2006-01-02", tradedate); err == nil {
				tradedate = t.Format("01/02/2006")
				line[1] = tradedate
			}

			exchange := line[5]
			line[5] = strings.ToLower(exchange)

			contract := line[6]
			line[6] = "c" + contract

			price, _ := strconv.ParseFloat(line[10], 64)
			line[10] = fmt.Sprintf("%.2f", price)

			unrealisedPL, _ := strconv.ParseFloat(strings.Replace(line[13], ",", "", -1), 64)
			line[13] = fmt.Sprintf("%.0f", unrealisedPL)

			commodity := line[17]
			line[17] = strings.ToLower(commodity)

			rows = append(rows, line)
		}
	}

	output.Write(destination+"/"+mainBillFiles[0], rows)
}

func mergeTrades(mainBillPath, subBillPath, destination string) {
	mainBillFiles, subBillFiles := getBillsInfo("Trade", mainBillPath, subBillPath)

	if len(mainBillFiles) <= 0 {
		fmt.Println("Trade 文件为空")
		return
	}

	if len(subBillFiles) <= 0 {
		fmt.Println("Trade 文件为空")
		return
	}

	temp, _ := readCSV(subBillPath + "/" + subBillFiles[0])
	header := temp[0]

	rows := [][]string{header}

	for _, filename := range subBillFiles {
		lines, err := readCSV(subBillPath + "/" + filename)
		if err != nil {
			fmt.Println("合并 Trade 错误")
			fmt.Println(err)
			break
		}

		for _, line := range lines[1:] {
			tradedate := line[1]
			if t, err := time.Parse("2006-01-02", tradedate); err == nil {
				tradedate = t.Format("01/02/2006")
				line[1] = tradedate
			}

			exchange := line[5]
			line[5] = strings.ToLower(exchange)

			contract := line[6]
			line[6] = "c" + contract

			buySell := line[15]
			if strings.ToLower(buySell) == "sale" {
				line[15] = "0"
			} else if strings.ToLower(buySell) == "buy" {
				line[15] = "1"
			}

			commodity := line[17]
			line[17] = strings.ToLower(commodity)

			rows = append(rows, line)
		}
	}

	output.Write(destination+"/"+mainBillFiles[0], rows)
}
