package cmd

import (
	"errors"
	"fmt"
	"github.com/pete911/awf/cmd/flag"
	"github.com/pete911/awf/internal/out"
	"github.com/pete911/awf/internal/store"
	"github.com/spf13/cobra"
	"net/netip"
	"os"
	"strings"
)

var (
	GlobalFlags flag.Global
	Root        = &cobra.Command{}
	Version     string
)

func init() {
	flag.InitPersistentFlags(Root, &GlobalFlags)
}

func NewTable() out.Table {
	return out.NewTable(os.Stdout, GlobalFlags.Trim)
}

func LoadFileStore() store.File {
	fileStorage, err := store.LoadFile()
	var notFound *store.NotFoundError
	if err != nil {
		if errors.As(err, &notFound) {
			fmt.Println("storage not initialized, run 'awf import' first")
			os.Exit(1)
		}
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return fileStorage
}

func IsVpcId(in string) bool {
	return strings.HasPrefix(in, "vpc-")
}

func IsSubnetId(in string) bool {
	return strings.HasPrefix(in, "subnet-")
}

func IsEniId(in string) bool {
	return strings.HasPrefix(in, "eni-")
}

func IsCIDR(in string) bool {
	if _, err := netip.ParsePrefix(in); err != nil {
		return false
	}
	return true
}

func IsIP(in string) bool {
	if _, err := netip.ParseAddr(in); err != nil {
		return false
	}
	return true
}
