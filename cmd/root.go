package cmd

import (
	"fmt"

	"github.com/bryantaolong/pdfx/cmd/commands"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pdfx",
	Short: "PDFX empowers you to merge, split and extract PDF files right in your terminal.",
}

var showVersion bool

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Print the version of pdfx")
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if showVersion {
			fmt.Println(GetVersion())
		}
		return nil
	}
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if showVersion {
			return nil
		}
		return cmd.Help()
	}

	rootCmd.AddCommand(commands.NewCmdMerge())
	rootCmd.AddCommand(commands.NewCmdSplit())
	rootCmd.AddCommand(commands.NewCmdExtract())
	rootCmd.AddCommand(NewCmdVersion())
}
