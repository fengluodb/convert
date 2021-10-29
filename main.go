package main

import (
	"convert/rule"
	"convert/txt"
	"fmt"
	"os"
	"regexp"
)

var r1, _ = regexp.Compile("./convert (.*txt) into epub")
var r2, _ = regexp.Compile("./convert (.*txt .*jpg) into epub")

func main() {
	command := ""
	for _, v := range os.Args {
		command += v + " "
	}
	fmt.Println(command)

	if r1.Match([]byte(command)) {
		file := txt.NewTxt(os.Args[1])
		rule.ConvertTxtToEpub(file, nil)
	}
	if r2.Match([]byte(command)) {
		file := txt.NewTxt(os.Args[1])
		cover, _ := os.Open(os.Args[2])
		rule.ConvertTxtToEpub(file, cover)
	}
}
