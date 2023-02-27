package image

type ImageInterface interface {
	ToMiddleImage() *MiddleImage
}

var ImageMap = map[string]func(path []string) (ImageInterface, error){
	"png": NewPNGs,
	"jpg": NewJPEGs,
	"gif": NewGIFs,
	"bmp": NewBMPs,
}
