package main

import (
	"convert/epub"
	"convert/rule"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("请输入更多参数")
		return
	}
	for i := 1; i < len(os.Args); i++ {
		if rule.IsTxt(os.Args[i]) {
			file, err := os.Open(os.Args[i])
			if err != nil {
				fmt.Println(err)
			}
			titleRule := rule.TxtTitleRegexp()
			epub.ConvertEpub(file, strings.Split(os.Args[i], ".")[0], titleRule, nil)
		}
	}
}
