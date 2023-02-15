package text

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/fengluodb/convert/utils"

	"github.com/pkg/errors"
)

var TxtTitleRegexp1, _ = regexp.Compile("^第[一二三四五六七八九零〇一二三四五六七八九十百千万a-zA-Z0-9]{1,7}[章节卷集部篇回].*|^前传[.]{0,10}|^正文[.]{0,10}")

type Txt struct {
	BookTitle              string
	ChapterTitles          []string
	ChapterTitleAndContent map[string]string
}

func NewTxt(path string) (*Txt, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.WithMessagef(err, "can't open %s", path)
	}

	defer file.Close()

	filepath := utils.GetFileNameFromPath(path)
	txt := &Txt{
		BookTitle:              strings.Split(filepath, ".")[0],
		ChapterTitles:          []string{},
		ChapterTitleAndContent: map[string]string{},
	}

	if err := txt.ParseTxt(bufio.NewReader(file)); err != nil {
		return nil, errors.WithMessagef(err, "can't parse %s", path)
	}

	return txt, nil
}

func (t *Txt) ParseTxt(buf *bufio.Reader) error {
	for {
		// 逐行读取，直到找到标题
		curline, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return errors.New("can't find chapter titles")
			} else {
				return err
			}
		}

		curline = strings.TrimSpace(curline)
		// fint first chapter title
		if TxtTitleRegexp1.Match([]byte(curline)) {
			if err := t.fillChapter(curline, buf); err != nil {
				return err
			}
			break
		}
	}

	return nil
}

func (t *Txt) fillChapter(firstTitle string, buf *bufio.Reader) error {
	title := firstTitle
	t.ChapterTitles = append(t.ChapterTitles, title)
	t.ChapterTitleAndContent[title] = ""

	for {
		curline, err := buf.ReadString('\n')
		curline = strings.TrimSpace(curline)
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}
		if TxtTitleRegexp1.Match([]byte(curline)) {
			title = curline
			t.ChapterTitles = append(t.ChapterTitles, title)
			t.ChapterTitleAndContent[title] = ""
		} else {
			// add \n，keep original format
			t.ChapterTitleAndContent[title] += curline + "\n"
		}
	}

	return nil
}

func (t *Txt) ToMiddleText() *MiddleText {
	return &MiddleText{
		BookTitle:              t.BookTitle,
		ChapterTitles:          t.ChapterTitles,
		ChapterTitleAndContent: t.ChapterTitleAndContent,
	}
}
