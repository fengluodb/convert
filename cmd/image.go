package cmd

import (
	"convert/image"
	"convert/utils"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	imageCmd = &cobra.Command{
		Use:   "image",
		Short: "image subcommand converts image file into another image format",
		Long:  "image subcommand converts image file into another image format",

		Run: func(cmd *cobra.Command, args []string) {
			if len(source) == 0 {
				if srcDir != "" && srcFormat != "" {
					source = utils.GetFilesWithSameSuffix(srcDir, srcFormat)
				}
				if len(source) == 0 {
					fmt.Println("source can't be empty")
					return
				}
			}

			originalFormat := strings.Split(source[0], ".")[1]

			original, err := image.ImageMap[originalFormat](source)
			if err != nil {
				fmt.Println("err: ", err)
			}
			middle := original.ToMiddleImage()
			if err := middle.To(format, output); err != nil {
				fmt.Println("err: ", err)
			}
		},
	}
)

func init() {
	imageCmd.Flags().StringSliceVarP(&source, "source", "s", nil, "the list of source files")
	imageCmd.Flags().StringVarP(&format, "format", "f", ".", "the format of output files")
	imageCmd.Flags().StringVarP(&output, "output", "o", ".", "the output dir")
	imageCmd.Flags().StringVarP(&srcDir, "srcdir", "", "", "the dir of output files, it doesn't work if source paramter has values")
	imageCmd.Flags().StringVarP(&srcFormat, "srcformat", "", "", "the format of src files, it works with srcdir paramter")

	imageCmd.MarkFlagRequired("format")

	rootCmd.AddCommand(imageCmd)
}
