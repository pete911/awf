package flag

import "github.com/spf13/cobra"

type Global struct {
	Trim bool
}

func InitPersistentFlags(cmd *cobra.Command, flags *Global) {
	cmd.PersistentFlags().BoolVar(
		&flags.Trim,
		"trim",
		true,
		"trim output",
	)
}
