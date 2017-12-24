package main

import "github.com/fengdu/billconverter/worker"

var src = "./src"
var destination = "./dst"

func main() {
	worker.Start(src, destination)
}
