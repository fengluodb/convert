package text

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

type Epub struct {
	BookTitle string
	RootDir   string
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

func NewEpub(output string, filename string) *Epub {
	rootDir := path.Join(output, filename)
	return &Epub{
		BookTitle: filename,
		RootDir:   rootDir,
		TextDir:   path.Join(rootDir, "text"),
		ImagesDir: path.Join(rootDir, "images"),
		StylesDir: path.Join(rootDir, "styles"),
		METADir:   path.Join(rootDir, "META-INF"),
		Opf: &OPF{
			BookTitle: filename,
			Date:      time.Now().Format("2006-01-02"),
		},
		Ncx: &NCX{
			BookTitle: filename,
		},
	}
}

// 初始化Epub的文件结构
// 我这里的结构并不是官方的标准结构，而是参考的其他文档建立的结构
// 详情看readme
func (e *Epub) InitializeEpub() {
	os.Mkdir(e.RootDir, os.ModePerm)
	os.Mkdir(e.TextDir, os.ModePerm)
	os.Mkdir(e.ImagesDir, os.ModePerm)
	os.Mkdir(e.StylesDir, os.ModePerm)
	os.Mkdir(e.METADir, os.ModePerm)
}

// 向text文件夹中添加文件
func (e *Epub) AddTextHtml(fileName, title, content string) {
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

// 向根目录中添加mimetype文件
func (e *Epub) AddMimetype() {

	file, _ := os.Create(path.Join(e.RootDir, "mimetype"))
	file.Write([]byte("application/epub+zip"))
	file.Close()
}

// 向META-INF目录中添加container.xml文件
func (e *Epub) AddContainerXml() {
	data, _ := tmpl.ReadFile("resources/container.xml")

	file, _ := os.Create(path.Join(e.METADir, "container.xml"))
	file.Write(data)
	file.Close()
}

// 向根目录中添加content.opf文件
func (e *Epub) AddContentOpf() {
	data, _ := tmpl.ReadFile("resources/content.opf")
	t, _ := template.New("opf").Parse(string(data))

	file, _ := os.Create(path.Join(e.RootDir, "content.opf"))
	t.Execute(file, e.Opf)
	file.Close()
}

// 向根目录中添加toc.ncx文件
func (e *Epub) AddTocNcx() {
	data, _ := tmpl.ReadFile("resources/toc.ncx")
	t, _ := template.New("ncx").Parse(string(data))

	file, _ := os.Create(path.Join(e.RootDir, "toc.ncx"))
	t.Execute(file, e.Ncx)
	file.Close()
}

// 向styles中添加css文件（默认的格式）
func (e *Epub) AddStyle() {
	data, _ := tmpl.ReadFile("resources/stylesheet.css")

	file, _ := os.Create(path.Join(e.StylesDir, "stylesheet.css"))
	file.Write(data)
	file.Close()
}

// 为epub添加缩略图
func (e *Epub) AddCover(src io.Reader) {
	file, _ := os.Create(path.Join(e.ImagesDir, "cover.jpg"))
	io.Copy(file, src)
	file.Close()
}

func (e *Epub) Zip() {
	fileName := e.RootDir + ".epub"

	// 预防：旧文件无法覆盖
	os.RemoveAll(fileName)

	// 创建：zip文件
	zipfile, _ := os.Create(fileName)
	defer zipfile.Close()

	// 打开：zip文件
	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	// 遍历路径信息
	filepath.Walk(e.RootDir, func(path string, info os.FileInfo, _ error) error {

		// 如果是源路径，提前进行下一个遍历
		if path == e.RootDir {
			return nil
		}

		// 获取：文件头信息
		header, _ := zip.FileInfoHeader(info)
		header.Name = strings.TrimPrefix(path, e.RootDir+`/`)

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
	os.RemoveAll(e.RootDir)
}
