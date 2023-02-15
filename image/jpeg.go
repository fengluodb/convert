package image

import (
	"convert/utils"
	"image"
	"image/jpeg"
	"os"

	"github.com/pkg/errors"
)

type JPEGs struct {
	filenames []string
	data      []image.Image
}

func NewJPEGs(paths []string) (ImageInterface, error) {
	data := []image.Image{}
	filenames := []string{}
	for _, path := range paths {
		filename := utils.GetFileNameFromPath(path)
		filenames = append(filenames, filename)

		file, err := os.Open(path)
		if err != nil {
			return nil, errors.WithMessagef(err, "can't open %s", path)
		}

		defer file.Close()
		image, err := jpeg.Decode(file)
		if err != nil {
			return nil, errors.WithMessagef(err, "can't decode %s", path)
		}
		data = append(data, image)
	}

	return &JPEGs{
		filenames: filenames,
		data:      data,
	}, nil
}

func (p *JPEGs) ToMiddleImage() *MiddleImage {
	return &MiddleImage{
		OriginalFormat:    "jpg",
		OriginalFilenames: p.filenames,
		ImageData:         p.data,
	}
}
