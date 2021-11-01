package rule

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

func ConvertImage(dst, src string) {
	srcfile, err := os.Open(src)
	if err != nil {
		fmt.Printf("打开%s失败", src)
		return
	}
	defer srcfile.Close()

	var originalImage image.Image

	imageStyle := strings.Split(src, ".")[1]
	switch imageStyle {
	case ".jpg":
		originalImage, err = jpeg.Decode(srcfile)
	case ".png":
		originalImage, err = png.Decode(srcfile)
	default:
		// 默认按jpg格式处理
		originalImage, err = jpeg.Decode(srcfile)
	}
	if err != nil {
		fmt.Println("图片解码时发送错误:", err)
		return
	}

	dstfile, err := os.Create(dst)
	if err != nil {
		fmt.Println("创建文件时发送错误:", err)
		return
	}
	defer dstfile.Close()
	err = jpeg.Encode(dstfile, originalImage, nil)
	if err != nil {
		fmt.Println("编码为jgp格式时发送错误:", err)
		return
	}
	fmt.Printf("%s转换%s成功", src, dst)
}
