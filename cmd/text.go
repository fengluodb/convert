package cmd

import (
	"convert/text"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	target string
	output string

	textCmd = &cobra.Command{
		Use:   "text",
		Short: "text subcommand converts text file into another format",
		Long:  "text subcommand converts text file into another format",

		Run: func(cmd *cobra.Command, args []string) {
			originalFormat := strings.Split(target, ".")[1]
			outputFormat := strings.Split(output, ".")[1]

			var middleText *text.MiddleText

			switch originalFormat {
			case "txt":
				txtfile, err := text.NewTxt(target)
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
	textCmd.Flags().StringVarP(&target, "target", "t", "", "target file")
	textCmd.Flags().StringVarP(&output, "output", "o", "", "output file")

	textCmd.MarkFlagRequired("target")
	textCmd.MarkFlagRequired("output")

	rootCmd.AddCommand(textCmd)
}
