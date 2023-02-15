package text

import (
	"fmt"
)

type MiddleText struct {
	BookTitle              string
	ChapterTitles          []string
	ChapterTitleAndContent map[string]string
}

func (m *MiddleText) ToEpub(output string, filename string) {
	dst := NewEpub(output, filename)
	dst.InitializeEpub()

	for i := 0; i < len(m.ChapterTitles); i++ {
		fileName := fmt.Sprintf("text%d.html", i)
		dst.AddTextHtml(fileName, m.ChapterTitles[i], m.ChapterTitleAndContent[m.ChapterTitles[i]])
		dst.Opf.Items += fmt.Sprintf(`<item href="text/text%d.html" id="id_%d" media-type="application/xhtml+xml"/>
			`, i, i)
		dst.Opf.Itemrefs += fmt.Sprintf(`<itemref idref="id_%d"/>
			`, i)
		dst.Ncx.NavPoints += fmt.Sprintf(`<navPoint id="id%d" playOrder="%d">
		   <navLabel>
		     <text>%s</text>
		   </navLabel>
		   <content src="text/text%d.html"/>
		</navPoint>`, i, i, m.ChapterTitles[i], i)
	}

	dst.AddMimetype()
	dst.AddContainerXml()
	dst.AddContentOpf()
	dst.AddTocNcx()
	dst.AddStyle()

	dst.Zip()
}
