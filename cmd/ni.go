package cmd

import (
	"fmt"
	"github.com/pete911/awf/internal/out"
	"github.com/pete911/awf/internal/types"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	niCmd = &cobra.Command{
		Use:   "ni",
		Short: "find network interface by IP or CIDR",
		Long:  "",
		Run:   runNi,
	}
)

func init() {
	Root.AddCommand(niCmd)
}

func runNi(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("no IP provided")
		return
	}

	fileStore := LoadFileStore()
	networkInterfaces, err := fileStore.DescribeNetworkInterfaces()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var matched types.NetworkInterfaces
	for _, arg := range args {
		if nis := findNetworkInterfaces(arg, networkInterfaces); len(nis) > 0 {
			matched = append(matched, nis...)
		}
	}

	if len(matched) == 0 {
		fmt.Printf("searched %d network interfaces, but none matched\n", len(networkInterfaces))
		return
	}

	printIp(matched)
}

func printIp(in types.NetworkInterfaces) {
	table := out.NewTable(os.Stdout)
	tableHeader := []string{"ACCOUNT ID", "AWS PROFILE", "ENI", "TYPE", "DESCRIPTION", "PRIVATE IP", "PUBLIC IP", "VPC", "SUBNET"}
	table.AddRow(tableHeader...)
	for _, v := range in {
		table.AddRow(
			v.Account.Id,
			v.Account.Profile,
			v.NetworkInterfaceId,
			v.Type,
			out.TrimTo(v.Description, 40),
			strings.Join(v.PrivateIpAddresses, ", "),
			v.PublicIP,
			v.VpcId,
			v.SubnetId,
		)
	}
	table.Print()
}

func findNetworkInterfaces(arg string, nis types.NetworkInterfaces) []types.NetworkInterface {
	if IsIP(arg) {
		return nis.GetByIp(arg)
	}
	if IsEniId(arg) {
		return nis.GetById(arg)
	}
	if IsCIDR(arg) {
		return nis.GetByCidr(arg)
	}
	fmt.Printf("argument %s can only be IP, CIDR, or network interface id", arg)
	os.Exit(1)
	return nil
}
