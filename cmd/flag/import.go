package flag

import "github.com/spf13/cobra"

type Import struct {
	Region string
}

func InitImportFlags(cmd *cobra.Command, flags *Import) {
	cmd.Flags().StringVar(
		&flags.Region,
		"region",
		"",
		"aws region",
	)
}
