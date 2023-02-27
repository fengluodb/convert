package image

import (
	"image"
	"os"

	"github.com/fengluodb/convert/utils"

	"github.com/pkg/errors"
	"golang.org/x/image/bmp"
)

type BMPs struct {
	filenames []string
	data      []image.Image
}

func NewBMPs(paths []string) (ImageInterface, error) {
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
		img, err := bmp.Decode(file)
		if err != nil {
			return nil, errors.WithMessagef(err, "can't decode %s", path)
		}
		data = append(data, img)
	}

	return &BMPs{
		filenames: filenames,
		data:      data,
	}, nil
}

func (p *BMPs) ToMiddleImage() *MiddleImage {
	return &MiddleImage{
		OriginalFormat:    "bmp",
		OriginalFilenames: p.filenames,
		ImageData:         p.data,
	}
}
