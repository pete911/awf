package cmd

import (
	"errors"
	"fmt"
	"github.com/pete911/awf/internal/store"
	"github.com/spf13/cobra"
	"net/netip"
	"os"
	"strings"
)

var (
	Root    = &cobra.Command{}
	Version string
)

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
