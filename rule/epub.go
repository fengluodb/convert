// 该文件中包含将其他格式转换为epub格式的规则
package rule

import (
	"convert/epub"
	"convert/txt"
	"fmt"
	"io"
)

// 将txt格式的网络小说转换为epub格式
func ConvertTxtToEpub(src *txt.TxtReader, cover io.ReadCloser) {
	defer cover.Close()

	src.ParseTxt()

	dst := epub.NewEpub(src.BookTitle)
	dst.InitializeEpub()

	for i := 0; i < len(src.ChapterTitles); i++ {
		fileName := fmt.Sprintf("text%d.html", i)
		dst.AddTextHtml(fileName, src.ChapterTitles[i], src.Chapter[src.ChapterTitles[i]])
		dst.Opf.Items += fmt.Sprintf(`<item href="text/text%d.html" id="id_%d" media-type="application/xhtml+xml"/>
			`, i, i)
		dst.Opf.Itemrefs += fmt.Sprintf(`<itemref idref="id_%d"/>
			`, i)
		dst.Ncx.NavPoints += fmt.Sprintf(`<navPoint id="id%d" playOrder="%d">
		   <navLabel>
		     <text>%s</text>
		   </navLabel>
		   <content src="text/text%d.html"/>
		</navPoint>`, i, i, src.ChapterTitles[i], i)
	}

	dst.AddMimetype()
	dst.AddContainerXml()
	dst.AddContentOpf()
	dst.AddTocNcx()
	dst.AddStyle()

	if cover != nil {
		dst.AddCover(cover)
	}

	dst.Zip()
}
