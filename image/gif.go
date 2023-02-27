package image

import (
	"fmt"
	"image"
	"image/gif"
	"os"

	"github.com/fengluodb/convert/utils"
)

type Gif struct {
	filenames []string
	data      []image.Image
}

func NewGIFs(paths []string) (ImageInterface, error) {
	filenames := []string{}
	data := []image.Image{}

	for _, path := range paths {
		gifName := utils.GetFileNameFromPath(path)

		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		inGif, err := gif.DecodeAll(file)
		if err != nil {
			return nil, err
		}
		for i := 0; i < len(inGif.Image); i++ {
			filename := fmt.Sprintf("%s-%d", gifName, i)
			filenames = append(filenames, filename)
			data = append(data, inGif.Image[i])
		}
	}

	return &Gif{
		filenames: filenames,
		data:      data,
	}, nil
}

func (g *Gif) ToMiddleImage() *MiddleImage {
	return &MiddleImage{
		OriginalFormat:    "gif",
		OriginalFilenames: g.filenames,
		ImageData:         g.data,
	}
}
