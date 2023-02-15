package cmd

import (
	"convert/text"
	"convert/utils"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	textCmd = &cobra.Command{
		Use:   "text",
		Short: "text subcommand converts text file into another format",
		Long:  "text subcommand converts text file into another format",

		// TODO: refine the code make it like image command
		Run: func(cmd *cobra.Command, args []string) {
			if output != "" {
				os.MkdirAll(output, 0777)
			}

			for _, path := range source {
				filename := utils.GetFileNameFromPath(path)
				originalFormat := strings.Split(path, ".")[1]

				var middleText *text.MiddleText

				switch originalFormat {
				case "txt":
					txtfile, err := text.NewTxt(path)
					if err != nil {
						panic(err)
					}
					middleText = txtfile.ToMiddleText()
				default:
					panic(fmt.Sprintf("can't read the %s filename", originalFormat))
				}

				switch format {
				case "epub":
					middleText.ToEpub(output, strings.Split(filename, ".")[0])
				default:
					panic(fmt.Sprintf("can't convert %s into %s", originalFormat, format))
				}
			}
		},
	}
)

func init() {
	textCmd.Flags().StringArrayVarP(&source, "source", "s", nil, "source files")
	textCmd.Flags().StringVarP(&format, "format", "f", "", "the target format")
	textCmd.Flags().StringVarP(&output, "output", "o", "", "the output dir")

	textCmd.MarkFlagRequired("source")
	textCmd.MarkFlagRequired("format")

	rootCmd.AddCommand(textCmd)
}
