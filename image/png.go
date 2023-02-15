package image

import (
	"image"
	"os"

	"github.com/fengluodb/convert/utils"

	"github.com/pkg/errors"
)

type PNGs struct {
	filenames []string
	data      []image.Image
}

func NewPNGs(paths []string) (ImageInterface, error) {
	filenames := []string{}
	data := []image.Image{}
	for _, path := range paths {
		filename := utils.GetFileNameFromPath(path)
		filenames = append(filenames, filename)

		file, err := os.Open(path)
		if err != nil {
			return nil, errors.WithMessagef(err, "can't open %s", path)
		}

		defer file.Close()
		img, _, err := image.Decode(file)
		if err != nil {
			return nil, errors.WithMessagef(err, "can't decode %s", path)
		}
		data = append(data, img)
	}

	return &PNGs{
		filenames: filenames,
		data:      data,
	}, nil
}

func (p *PNGs) ToMiddleImage() *MiddleImage {
	return &MiddleImage{
		OriginalFormat:    "png",
		OriginalFilenames: p.filenames,
		ImageData:         p.data,
	}
}
