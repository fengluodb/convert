package image

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path"

	"github.com/disintegration/imaging"
	"github.com/pkg/errors"
	"golang.org/x/image/bmp"
)

type MiddleImage struct {
	OriginalFormat    string
	OriginalFilenames []string
	ImageData         []image.Image
}

func (m *MiddleImage) To(format string, outputDir string) error {
	if outputDir != "." {
		os.MkdirAll(outputDir, 0777)
	}

	switch {
	case format == "png":
		return m.ToPng(outputDir)
	case format == "jpg" || format == "jpeg":
		return m.ToJPEG(outputDir)
	case format == "bmp":
		return m.ToBMP(outputDir)
	case format == "gif":
		return m.ToGif(outputDir)
	}
	return fmt.Errorf("can't convet %s into %s", m.OriginalFormat, format)
}

func (m *MiddleImage) ToPng(outputDir string) error {
	return ToPng(outputDir, m.OriginalFilenames, m.ImageData)
}

func (m *MiddleImage) ToJPEG(outputDir string) error {
	return ToJPEG(outputDir, m.OriginalFilenames, m.ImageData)
}

func (m *MiddleImage) ToBMP(outputDir string) error {
	return ToBMP(outputDir, m.OriginalFilenames, m.ImageData)
}

func (m *MiddleImage) ToGif(output string) error {
	return ToGif(output, m.ImageData)
}

func (m *MiddleImage) Resize(width, height int) {
	if m.OriginalFormat == "gif" {
		fmt.Println("current don't support resize gif")
	}

	for i, v := range m.ImageData {
		resize_height := height
		resize_width := width
		if resize_height == 0 {
			resize_height = v.Bounds().Dx()
		}
		if resize_width == 0 {
			resize_width = v.Bounds().Dy()
		}
		m.ImageData[i] = imaging.Resize(v, resize_width, resize_height, imaging.Lanczos)
	}
}

func ToPng(outputDir string, filenames []string, imageData []image.Image) error {
	for i, img := range imageData {
		filepath := path.Join(outputDir, filenames[i]+".png")
		file, err := os.Create(filepath)
		if err != nil {
			return errors.WithMessagef(err, "can't create %s", filepath)
		}

		defer file.Close()
		err = png.Encode(file, img)
		if err != nil {
			return errors.WithMessagef(err, "failed to encode into %s", filepath)
		}
	}
	return nil
}

func ToJPEG(outputDir string, filenames []string, imageData []image.Image) error {
	for i, img := range imageData {
		filepath := path.Join(outputDir, filenames[i]+".jpg")
		file, err := os.Create(filepath)
		if err != nil {
			return errors.WithMessagef(err, "can't create %s", filepath)
		}

		defer file.Close()
		err = jpeg.Encode(file, img, nil)
		if err != nil {
			return errors.WithMessagef(err, "failed to encode into %s", filepath)
		}
	}
	return nil
}

func ToBMP(outputDir string, filenames []string, imageData []image.Image) error {
	for i, img := range imageData {
		filepath := path.Join(outputDir, filenames[i]+".bmp")
		file, err := os.Create(filepath)
		if err != nil {
			return errors.WithMessagef(err, "can't create %s", filepath)
		}

		defer file.Close()
		err = bmp.Encode(file, img)
		if err != nil {
			return errors.WithMessagef(err, "failed to encode into %s", filepath)
		}
	}
	return nil
}

func ToGif(outputDir string, imageData []image.Image) error {
	g := &gif.GIF{
		LoopCount: len(imageData),
	}

	for _, img := range imageData {
		g.Image = append(g.Image, img.(*image.Paletted))
		g.Delay = append(g.Delay, 0)
	}

	path := path.Join(outputDir, "output.gif")
	f, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	if err := gif.EncodeAll(f, g); err != nil {
		return errors.WithMessagef(err, "can't encode into %s", path)
	}
	return nil
}
