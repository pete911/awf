package cmd

import (
	"github.com/pete911/awf/cmd/flag"
	"github.com/spf13/cobra"
)

var (
	GlobalFlags flag.Global
	Root        = &cobra.Command{}
	Version     string
)

func init() {
	flag.InitPersistentFlags(Root, &GlobalFlags)
}
