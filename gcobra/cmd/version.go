package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var verbose *bool

func init() {
	verbose = versionCmd.PersistentFlags().BoolP("version", "v", false, "verbose output")

	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Args: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("%+v\n", args)
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
		fmt.Println(*verbose)
		fmt.Println(args)
	},
}
