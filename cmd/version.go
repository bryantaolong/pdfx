package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionStr = "v0.1.1"

func GetVersion() string {
	return versionStr
}

func NewCmdVersion() *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Short:   "Print the version of pdfx",
		Aliases: []string{"v"},
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(versionStr)
		},
	}
}
