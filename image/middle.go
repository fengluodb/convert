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
)

type MiddleImage struct {
	OriginalFormat    string
	OriginalFilenames []string

	ImageData []image.Image
	GifData   []*gif.GIF
}

func (m *MiddleImage) To(format string, outputDir string) error {
	if outputDir != "." {
		os.MkdirAll(outputDir, 0777)
	}

	switch {
	case format == "png":
		if m.OriginalFormat == "gif" {
			return m.GifToImage(format, outputDir)
		}
		return m.ToPng(outputDir)
	case format == "jpg" || format == "jpeg":
		if m.OriginalFormat == "gif" {
			return m.GifToImage(format, outputDir)
		}
		return m.ToJPEG(outputDir)
	case format == "gif":
		if m.OriginalFormat == "gif" {
			for _, gif := range m.GifData {
				m.ImageData = append(m.ImageData, GifToImage(gif)...)
			}
		}
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

func (m *MiddleImage) ToGif(output string) error {
	return ToGif(output, m.OriginalFilenames[0], m.ImageData)
}

func (m *MiddleImage) GifToImage(format string, outputDir string) error {
	var fn func(outputDir string, filenames []string, imageData []image.Image) error
	if format == "png" {
		fn = ToPng
	} else if format == "jpg" || format == "jpeg" {
		fn = ToJPEG
	} else {
		return fmt.Errorf("don't support convert git into %s", format)
	}

	for i, inGif := range m.GifData {
		imageData := GifToImage(inGif)
		filenames := []string{}
		for j := 0; j < len(imageData); j++ {
			filename := fmt.Sprintf("%s-%d", m.OriginalFilenames[i], j)
			filenames = append(filenames, filename)
		}
		if err := fn(outputDir, filenames, imageData); err != nil {
			return err
		}
	}
	return nil
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

func GifToImage(inGif *gif.GIF) []image.Image {
	data := []image.Image{}
	for _, img := range inGif.Image {
		data = append(data, img)
	}

	return data
}

func ToGif(outputDir string, filename string, imageData []image.Image) error {
	g := &gif.GIF{
		LoopCount: len(imageData),
	}

	for _, img := range imageData {
		g.Image = append(g.Image, img.(*image.Paletted))
		g.Delay = append(g.Delay, 0)
	}

	path := path.Join(outputDir, filename+".gif")
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
