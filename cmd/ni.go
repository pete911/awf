package cmd

import (
	"fmt"
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
		fmt.Println("no argument provided")
		return
	}

	fileStore := LoadFileStore()
	vpcs, err := fileStore.DescribeVpcs()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	sunbets, err := fileStore.DescribeSubnets()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	nis, err := fileStore.DescribeNetworkInterfaces()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var matched types.NetworkInterfaces
	for _, arg := range args {
		if nis := findNetworkInterfaces(arg, nis); len(nis) > 0 {
			matched = append(matched, nis...)
		}
	}

	if len(matched) == 0 {
		fmt.Printf("searched %d network interfaces, but none matched\n", len(nis))
		return
	}

	printNi(matched, vpcs, sunbets)
}

func printNi(nis types.NetworkInterfaces, vpcs types.Vpcs, subnets types.Subnets) {
	table := NewTable()
	tableHeader := []string{"ACCOUNT ID", "AWS PROFILE", "ENI", "TYPE", "DESCRIPTION", "PRIVATE IP", "PUBLIC IP", "VPC ID", "VPC NAME", "SUBNET ID", "SUBNET NAME"}
	table.AddRow(tableHeader...)
	for _, v := range nis {
		var vpcName string
		if x := vpcs.GetById(v.VpcId); len(x) != 0 {
			vpcName = x[0].Name
		}
		var subnetName string
		if x := subnets.GetById(v.SubnetId); len(x) != 0 {
			subnetName = x[0].Name
		}

		table.AddRow(
			v.Account.Id,
			v.Account.Profile,
			v.NetworkInterfaceId,
			v.Type,
			v.Description,
			strings.Join(v.PrivateIpAddresses, ", "),
			v.PublicIP,
			v.VpcId,
			vpcName,
			v.SubnetId,
			subnetName,
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
