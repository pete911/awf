package cmd

import (
	"fmt"
	"github.com/pete911/awf/cmd/flag"
	"github.com/pete911/awf/internal/store"
	"github.com/spf13/cobra"
	"os"
)

var (
	importFlags flag.Import
	importCmd   = &cobra.Command{
		Use:   "import",
		Short: "import aws resources to local storage",
		Long:  "",
		Run:   runImport,
	}
)

func init() {
	flag.InitImportFlags(importCmd, &importFlags)
	Root.AddCommand(importCmd)
}

func runImport(cmd *cobra.Command, _ []string) {
	logger := GlobalFlags.Logger().With("cmd", cmd.Use)

	if err := store.Import(logger, importFlags.Region); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
