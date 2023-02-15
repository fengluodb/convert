package cmd

import (
	"convert/text"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	textCmd = &cobra.Command{
		Use:   "text",
		Short: "text subcommand converts text file into another format",
		Long:  "text subcommand converts text file into another format",

		// TODO: add support for current command style
		Run: func(cmd *cobra.Command, args []string) {
			originalFormat := strings.Split(source[0], ".")[1]
			outputFormat := strings.Split(output, ".")[1]

			var middleText *text.MiddleText

			switch originalFormat {
			case "txt":
				txtfile, err := text.NewTxt(source[0])
				if err != nil {
					panic(err)
				}
				middleText = txtfile.ToMiddleText()
			default:
				panic(fmt.Sprintf("can't read the %s filename", originalFormat))
			}

			switch outputFormat {
			case "epub":
				middleText.ToEpub(output)
			default:
				panic(fmt.Sprintf("can't convert %s into %s", originalFormat, outputFormat))
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
	textCmd.MarkFlagRequired("output")

	rootCmd.AddCommand(textCmd)
}
