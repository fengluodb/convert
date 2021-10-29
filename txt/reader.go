package txt

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

type TxtReader struct {
	BookTitle     string // 本书标题
	file          *os.File
	buf           *bufio.Reader
	ChapterTitles []string          // 章节标题
	Chapter       map[string]string // 标题：内容
}

var TxtTitleRegexp1, _ = regexp.Compile("^第[一二三四五六七八九零〇一二三四五六七八九十百千万a-zA-Z0-9]{1,7}[章节卷集部篇回].*|^前传[.]{0,10}|^正文[.]{0,10}")

func NewTxt(filepath string) *TxtReader {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("打开txt文件失败")
		return nil
	}

	return &TxtReader{
		BookTitle:     strings.Split(filepath, ".")[0],
		file:          file,
		buf:           bufio.NewReader(file),
		ChapterTitles: []string{},
		Chapter:       map[string]string{},
	}
}

func (t *TxtReader) ParseTxt() {
	for {
		// 逐行读取，直到找到标题
		curline, err := t.buf.ReadString('\n')
		curline = strings.TrimSpace(curline)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println("逐行读取Txt时发生错误")
				return
			}
		}

		// 找到第一个标题
		if TxtTitleRegexp1.Match([]byte(curline)) {
			t.fillChapter(curline)
			break
		}
	}
	t.file.Close()
}

func (t *TxtReader) fillChapter(firstTitle string) {
	title := firstTitle
	t.ChapterTitles = append(t.ChapterTitles, title)
	t.Chapter[title] = ""

	for {
		curline, err := t.buf.ReadString('\n')
		curline = strings.TrimSpace(curline)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println("读取章节内容时发生错误")
			}
		}
		if TxtTitleRegexp1.Match([]byte(curline)) {
			title = curline
			t.ChapterTitles = append(t.ChapterTitles, title)
			t.Chapter[title] = ""
		} else {
			// 重新加入换行符，保持源格式
			t.Chapter[title] += curline + "\n"
		}
	}
}
