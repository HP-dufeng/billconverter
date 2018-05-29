package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	mainflag := flag.String("main", "./dst", "主账单目录")
	subflag := flag.String("sub", "./dst_sub", "子帐单目录")
	destinationflag := flag.String("dst_zip", "./dst_zip", "zip文件目录")

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
	fmt.Println("zip文件目录: ", destination)

	// Read bills from folder
	var mainBillFiles, subBillFiles []os.FileInfo
	var err error
	if mainBillFiles, err = ioutil.ReadDir(mainBillPath); err != nil {
		fmt.Println("读取主账单目录错误", err)
		return
	}
	if subBillFiles, err = ioutil.ReadDir(subBillPath); err != nil {
		fmt.Println("读取子帐单目录错误", err)
		return
	}

	balances, balanceZipFileName := filterFiles("Balance", mainBillPath, subBillPath, destination, mainBillFiles, subBillFiles)
	poses, posZipFileName := filterFiles("Pos", mainBillPath, subBillPath, destination, mainBillFiles, subBillFiles)
	trades, tradeZipFileName := filterFiles("Trade", mainBillPath, subBillPath, destination, mainBillFiles, subBillFiles)

	ziper(balances, balanceZipFileName)
	ziper(poses, posZipFileName)
	ziper(trades, tradeZipFileName)
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
			return "ERROR: 创建压缩目录错误 ", err
		}
	}

	return "", nil
}

func filterFiles(pattern, mainFilePath, subFilePath, destinationPath string, mainBillFiles, subBillFiles []os.FileInfo) ([]string, string) {
	var files []string
	var zipFile string
	for _, f := range mainBillFiles {
		if f.IsDir() {
			continue
		}
		if strings.Contains(f.Name(), pattern) {
			files = append(files, mainFilePath+"/"+f.Name())
			zipFile = destinationPath + "/" + strings.Replace(f.Name(), ".csv", ".zip", -1)
		}
	}
	for _, f := range subBillFiles {
		if f.IsDir() {
			continue
		}
		if strings.Contains(f.Name(), pattern) {
			files = append(files, subFilePath+"/"+f.Name())
		}
	}

	return files, zipFile
}

func ziper(filePaths []string, dest string) error {
	d, _ := os.Create(dest)
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()

	for _, filepath := range filePaths {
		file, err := os.Open(filepath)
		defer file.Close()
		if err != nil {
			fmt.Println("读取错误: ", filepath, err)
			return err
		}

		err = compress(file, "", w)
		if err != nil {
			return err
		}

	}
	return nil
}

func compress(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	header.Name = prefix + "/" + header.Name
	if err != nil {
		return err
	}
	writer, err := zw.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, file)
	file.Close()
	if err != nil {
		return err
	}
	return nil
}
