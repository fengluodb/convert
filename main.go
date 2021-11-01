package main

import (
	"convert/rule"
	"convert/txt"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

var r1, _ = regexp.Compile("./convert (.*txt) into epub")
var r2, _ = regexp.Compile("./convert .*txt (.*jpg|.*png) into epub")
var r3, _ = regexp.Compile("./convert (.*jpg|.*png) into (.*jpg|.*png)")

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
		var cover *os.File
		if strings.Split(os.Args[2], ".")[1] == "png" {
			defer os.Remove("tmp.jpg")
			rule.ConvertImage("tmp.jpg", os.Args[2])
			time.Sleep(10 * time.Second)
			cover, _ = os.Open("tmp.jpg")
			file := txt.NewTxt(os.Args[1])
			rule.ConvertTxtToEpub(file, cover)
		} else {
			cover, _ = os.Open(os.Args[2])
			file := txt.NewTxt(os.Args[1])
			rule.ConvertTxtToEpub(file, cover)
		}
	}

	if r3.Match([]byte(command)) {
		rule.ConvertImage(os.Args[3], os.Args[1])
	}
}
