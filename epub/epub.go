package epub

import (
	"archive/zip"
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

func ConvertEpub(src *os.File, booktitle string, titleRegrex, contentRegrex *regexp.Regexp) {
	InitializeEpub(booktitle)
	ReadFromTxt(src, booktitle, titleRegrex, contentRegrex)
	Zip(booktitle, booktitle+".epub")
}

func InitializeEpub(bookTitle string) error {
	rootPath := path.Join(bookTitle)
	mimetypePath := path.Join(rootPath, "mimetype")
	metaPath := path.Join(rootPath, "META-INf")
	containerPath := path.Join(metaPath, "container.xml")
	oebpsPath := path.Join(rootPath, "OEBPS")

	// 创建根目录
	err := os.Mkdir(rootPath, os.ModePerm)
	if err != nil {
		return errors.New("failed to create root directory")
	}

	// 创建mimetype文件
	mimetype, err := os.Create(mimetypePath)
	defer func() {
		mimetype.Close()
	}()
	if err != nil {
		return errors.New("failed to create mimetype file")
	} else {
		_, err := mimetype.Write([]byte("application/equb+zip"))
		if err != nil {
			return errors.New("failed to write mimetype")
		}
	}

	// 创建META_INf文件夹
	err = os.Mkdir(metaPath, os.ModePerm)
	if err != nil {
		return errors.New("failed to create META_INF directory")
	}

	// 创建container.xml文件
	container, err := os.Create(containerPath)
	defer func() {
		container.Close()
	}()
	if err != nil {
		return errors.New("failed to create container.xml file")
	} else {
		_, err := container.Write([]byte(`<?xml version="1.0" encoding="UTf-8" ?>
<container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container">
<rootfiles> 
    <rootfile full-path="OEBPS/content.opf" media-type="application/oebps-package+xml"/> </rootfiles>
</container>`))
		if err != nil {
			return errors.New("failed to write container.xml")
		}
	}

	// 创建OEPBS文件夹
	err = os.Mkdir(oebpsPath, os.ModePerm)
	if err != nil {
		return errors.New("failed to create OEBPS directory")
	}

	return nil
}

var ncxFileTemplate = `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE ncx PUBLIC "-//NISO//DTD ncx 2005-1//EN" "http://www.daisy.org/z3986/2005/ncx-2005-1.dtd">
<ncx version="2005-1" xmlns="http://www.daisy.org/z3986/2005/ncx/">
<head>
  <meta name="dtb:uid" content=" "/>
  <meta name="dtb:depth" content="-1"/>
  <meta name="dtb:totalPageCount" content="0"/>
  <meta name="dtb:maxPageNumber" content="0"/>
</head>
 <docTitle><text>%s</text></docTitle>
<navMap>
%s
</navMap>
</ncx>`
var navTemplate = `<navPoint id="id%d" playOrder="%d">
       <navLabel>
         <text>%s</text>
       </navLabel>
       <content src="chapter%d.html"/>
    </navPoint>`

var chapterTemplate = `<html xmlns=http://www.w3.org/1999/xhtml xml:lang=zh-CN>
	<head>
        <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
        <link rel="stylesheet" type="text/css" href="css/main.css" />
        <title>%s</title>
    </head>
	<body>
		<div>
		%s
		</div>
	</body>
<html>`
var chapterContentTemplate = "\t<p>\n%s\n</p>\n"

func ReadFromTxt(src *os.File, booktitle string, titleRegrex, contentRegrex *regexp.Regexp) error {
	defer src.Close()
	buf := bufio.NewReader(src)

	if titleRegrex == nil {
		return errors.New("need titleRegrex")
	}

	oepbsDir := path.Join(booktitle, "OEBPS")
	ncxFile, err := os.Create(path.Join(oepbsDir, "toc.ncx"))
	navContent := ""
	defer func() {
		ncxFile.Write([]byte(fmt.Sprintf(ncxFileTemplate, booktitle, navContent)))
		ncxFile.Close()
	}()
	if err != nil {
		return errors.New("failed create toc.ncx")
	}

	curChapter := 0
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF {
				fmt.Println("file read ok!")
				break
			} else {
				fmt.Println()
				return errors.New("read file error")
			}
		}
		fmt.Println(line)
	label:
		curChapterContent := ""
		if titleRegrex.Match([]byte(line)) {
			navContent += fmt.Sprintf(navTemplate, curChapter, curChapter, line, curChapter)
			curChapterPath := path.Join(oepbsDir, fmt.Sprintf("chapter%d.html", curChapter))
			curChapterFile, err := os.Create(curChapterPath)
			if err != nil {
				return fmt.Errorf("fail create chapter%d.html", curChapter)
			}
			for {
				line, err = buf.ReadString('\n')
				line = strings.TrimSpace(line)
				if err != nil {
					if err == io.EOF {
						fmt.Println("File read ok!")
						break
					} else {
						fmt.Println("Read file error!", err)
						return err
					}
				}
				if titleRegrex.Match([]byte(line)) {
					break
				} else {
					curChapterContent += fmt.Sprintf(chapterContentTemplate, line)
				}
			}
			curChapterFile.Write([]byte(fmt.Sprintf(chapterTemplate, booktitle, curChapterContent)))
			curChapterFile.Close()
			curChapter++
			goto label
		}
	}

	opfPath := path.Join(oepbsDir, "content.opf")
	writeOpf(opfPath, booktitle, curChapter)
	return nil

}

var opfFileTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<package xmlns:opf="http://www.idpf.org/2007/opf" unique-identifier="bookid" xmlns="http://www.idpf.org/2007/opf" version="2.0">
   <metadata xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:opf="http://www.idpf.org/2007/opf" xmlns:dcterms="http://purl.org/dc/terms/" xmlns:dc="http://purl.org/dc/elements/1.1/">
        <dc:language>zh</dc:language>
        <dc:title>%s</dc:title>
   </metadata>
   <manifest>
		%s
   </manifest>
   <spine toc = "ncx">
   		%s
   </spine>
</package>`

func writeOpf(filepath, booktitle string, n int) {
	mainfest := ""
	toc := ""
	for i := 0; i <= n; i++ {
		mainfest += fmt.Sprintf(`<item href="chapter%d.html" id="id_%d" media-type="application/xhtml+xml"/>
		`, i, i)
		toc += fmt.Sprintf(`<itemref idref="id_%d"/>
		`, i)
	}
	mainfest += `<item href="toc.ncx" id="ncx" media-type="application/x-dtbncx+xml"/>
	`

	file, err := os.Create(filepath)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		file.Write([]byte(fmt.Sprintf(opfFileTemplate, booktitle, mainfest, toc)))
	}
	file.Close()
}

func Zip(srcDir, zipFileName string) {
	// 预防：旧文件无法覆盖
	os.RemoveAll(zipFileName)

	// 创建：zip文件
	zipfile, _ := os.Create(zipFileName)
	defer zipfile.Close()

	// 打开：zip文件
	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	// 遍历路径信息
	filepath.Walk(srcDir, func(path string, info os.FileInfo, _ error) error {

		// 如果是源路径，提前进行下一个遍历
		if path == srcDir {
			return nil
		}

		// 获取：文件头信息
		header, _ := zip.FileInfoHeader(info)
		header.Name = strings.TrimPrefix(path, srcDir+`\`)

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
	os.RemoveAll(srcDir)
}
