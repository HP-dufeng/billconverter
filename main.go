// Convert bill txt files into csv format
// The core logic is in this file: worker.go
// Logic Description:
// 		input ==> convert ==> output
// input: read bill file content from src folder
// convert: extract the segments of [Trade Confirmation], [Gathered Open Positions], [Financial Situation] ...
// output: write segments into csv file
package main

import (
	"flag"
	"fmt"

	"github.com/fengdu/billconverter/worker"
)

func main() {
	src := flag.String("src", "./src", "src folder")
	destination := flag.String("dst", "./dst", "dst folder")

	flag.Parse()
	fmt.Printf("Src folder: %s\n", *src)
	fmt.Printf("Destination folder: %s\n", *destination)

	worker.Start(*src, *destination)
}
