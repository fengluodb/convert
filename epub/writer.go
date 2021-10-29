package epub

import (
	"archive/zip"
	"bufio"
	"embed"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

type EpubWriter struct {
	BookTitle string
	TextDir   string
	ImagesDir string
	StylesDir string
	METADir   string

	Opf *OPF
	Ncx *NCX
}

type OPF struct {
	BookTitle string
	Date      string
	Items     string
	Itemrefs  string
}

type NCX struct {
	BookTitle string
	NavPoints string
}

type TEXT struct {
	BookTitle  string
	Title      string
	Paragraphs string
}

//go:embed resources
var tmpl embed.FS

func NewEpub(title string) *EpubWriter {
	return &EpubWriter{
		BookTitle: title,
		TextDir:   path.Join(title, "text"),
		ImagesDir: path.Join(title, "images"),
		StylesDir: path.Join(title, "styles"),
		METADir:   path.Join(title, "META-INF"),
		Opf: &OPF{
			BookTitle: title,
			Date:      time.Now().Format("2006-01-02"),
		},
		Ncx: &NCX{
			BookTitle: title,
		},
	}
}

// 初始化Epub的文件结构
// 我这里的结构并不是官方的标准结构，而是参考的其他文档建立的结构
// 详情看readme
func (e *EpubWriter) InitializeEpub() {
	os.Mkdir(e.BookTitle, os.ModePerm)
	os.Mkdir(e.TextDir, os.ModePerm)
	os.Mkdir(e.ImagesDir, os.ModePerm)
	os.Mkdir(e.StylesDir, os.ModePerm)
	os.Mkdir(e.METADir, os.ModePerm)
}

// 向Text文件夹中添加文件
func (e *EpubWriter) AddTextHtml(fileName, title, content string) {
	buf := bufio.NewReader(strings.NewReader(content))
	paragraphs := ""
	for {
		p, err := buf.ReadString('\n')
		p = strings.TrimSpace(p)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println("向text文件夹中添加文件失败")
				return
			}
		}
		paragraphs += fmt.Sprintf("<p>%s</p>", p)
	}
	text := TEXT{
		BookTitle:  e.BookTitle,
		Title:      title,
		Paragraphs: paragraphs,
	}

	data, _ := tmpl.ReadFile("resources/chapter.html")
	t, _ := template.New("text").Parse(string(data))

	file, _ := os.Create(path.Join(e.TextDir, fileName))
	t.Execute(file, text)
	file.Close()
}

func (e *EpubWriter) AddMimetype() {

	file, _ := os.Create(path.Join(e.BookTitle, "mimetype"))
	file.Write([]byte("application/epub+zip"))
	file.Close()
}

func (e *EpubWriter) AddContainerXml() {
	data, _ := tmpl.ReadFile("resources/container.xml")

	file, _ := os.Create(path.Join(e.METADir, "container.xml"))
	file.Write(data)
	file.Close()
}

func (e *EpubWriter) AddContentOpf() {
	data, _ := tmpl.ReadFile("resources/content.opf")
	t, _ := template.New("opf").Parse(string(data))

	file, _ := os.Create(path.Join(e.BookTitle, "content.opf"))
	t.Execute(file, e.Opf)
	file.Close()
}

func (e *EpubWriter) AddTocNcx() {
	data, _ := tmpl.ReadFile("resources/toc.ncx")
	t, _ := template.New("ncx").Parse(string(data))

	file, _ := os.Create(path.Join(e.BookTitle, "toc.ncx"))
	t.Execute(file, e.Ncx)
	file.Close()
}

func (e *EpubWriter) AddStyle() {
	data, _ := tmpl.ReadFile("resources/stylesheet.css")

	file, _ := os.Create(path.Join(e.StylesDir, "stylesheet.css"))
	file.Write(data)
	file.Close()
}

func (e *EpubWriter) AddCover(src io.Reader) {
	file, _ := os.Create(path.Join(e.ImagesDir, "cover.jpg"))
	io.Copy(file, src)
	file.Close()
}

func (e *EpubWriter) Zip() {
	fileName := e.BookTitle + ".epub"

	// 预防：旧文件无法覆盖
	os.RemoveAll(fileName)

	// 创建：zip文件
	zipfile, _ := os.Create(fileName)
	defer zipfile.Close()

	// 打开：zip文件
	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	// 遍历路径信息
	filepath.Walk(e.BookTitle, func(path string, info os.FileInfo, _ error) error {

		// 如果是源路径，提前进行下一个遍历
		if path == e.BookTitle {
			return nil
		}

		// 获取：文件头信息
		header, _ := zip.FileInfoHeader(info)
		header.Name = strings.TrimPrefix(path, e.BookTitle+`\`)

		// 判断：文件是不是文件夹
		if info.IsDir() {
			header.Name += `/`
		} else {
			// 设置：zip的文件压缩算法
			header.Method = zip.Deflate
		}

		// 创建：压缩包头部信息
		writer, _ := archive.CreateHeader(header)
		if !info.IsDir() {
			file, _ := os.Open(path)
			defer file.Close()
			io.Copy(writer, file)
		}
		return nil
	})
	os.RemoveAll(e.BookTitle)
}
