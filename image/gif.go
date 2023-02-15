package image

import (
	"convert/utils"
	"image/gif"
	"os"
)

type Gif struct {
	filenames []string
	data      []*gif.GIF
}

func NewGIFs(paths []string) (ImageInterface, error) {
	filenames := []string{}
	data := []*gif.GIF{}

	for _, path := range paths {
		filename := utils.GetFileNameFromPath(path)
		filenames = append(filenames, filename)

		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		inGif, err := gif.DecodeAll(file)
		if err != nil {
			return nil, err
		}
		data = append(data, inGif)
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
		GifData:           g.data,
	}
}
