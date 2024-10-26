package flag

import (
	"fmt"
	"github.com/pete911/awf/internal"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
)

type Global struct {
	logLevel string
}

func (f Global) Logger() *slog.Logger {
	l, err := internal.NewLogger(f.logLevel)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return l
}

func InitPersistentFlags(cmd *cobra.Command, flags *Global) {
	cmd.PersistentFlags().StringVar(
		&flags.logLevel,
		"log-level",
		"info",
		"log level - debug, info, warn, error",
	)
}
