package main

import (
	"convert/rule"
	"convert/txt"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Println("输入格式错误, 请参考标准输入:")
		fmt.Println("\t./convert a.txt或")
		fmt.Println("\t./convert a.txt b.png")
	}

	file := txt.NewTxt(os.Args[1])
	rule.ConvertTxtToEpub(file, nil)
}
