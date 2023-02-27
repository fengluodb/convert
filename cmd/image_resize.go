package cmd

import (
	"fmt"
	"strings"

	"github.com/fengluodb/convert/image"
	"github.com/spf13/cobra"
)

var imageResizeCmd = &cobra.Command{
	Use:   "resize",
	Short: "image resize subcommand resize the image to the desired",
	Long:  "image resize subcommand resize the image to the desired",

	Run: func(c *cobra.Command, args []string) {
		if imageHeight == 0 && imageWidth == 0 {
			fmt.Println("please give the height or width")
		}

		originalFormat := strings.Split(source[0], ".")[1]

		original, err := image.ImageMap[originalFormat](source)
		if err != nil {
			fmt.Println("err: ", err)
		}
		middle := original.ToMiddleImage()
		middle.Resize(imageWidth, imageHeight)
		if err := middle.To(originalFormat, output); err != nil {
			fmt.Println("err: ", err)
		}
	},
}

func init() {
	imageResizeCmd.Flags().StringSliceVarP(&source, "source", "s", nil, "the list of source files")
	imageResizeCmd.Flags().StringVarP(&output, "output", "o", ".", "the output dir")
	imageResizeCmd.Flags().IntVar(&imageHeight, "height", 0, "the height of the image to resize")
	imageResizeCmd.Flags().IntVar(&imageWidth, "width", 0, "the width of the image to resize")

	imageResizeCmd.MarkFlagRequired("source")

	imageCmd.AddCommand(imageResizeCmd)
}
