package cmd

import (
	"fmt"
	"github.com/pete911/awf/internal/out"
	"github.com/pete911/awf/internal/types"
	"github.com/spf13/cobra"
	"os"
)

var (
	subnetCmd = &cobra.Command{
		Use:   "subnet",
		Short: "find subnet by IP or CIDR",
		Long:  "",
		Run:   runSubnet,
	}
)

func init() {
	Root.AddCommand(subnetCmd)
}

func runSubnet(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("no argument provided")
		return
	}

	fileStore := LoadFileStore()
	accounts, err := fileStore.ListAccounts()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	vpcs, err := fileStore.DescribeVpcs()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	subnets, err := fileStore.DescribeSubnets()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	nis, err := fileStore.DescribeNetworkInterfaces()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var matched types.Subnets
	for _, arg := range args {
		if found := findSubnets(arg, subnets); len(found) > 0 {
			matched = append(matched, found...)
		}
	}

	if len(matched) == 0 {
		fmt.Printf("searched %d vpcs, but none matched\n", len(vpcs))
		return
	}

	printSubnets(nis, vpcs, matched, accounts)
}

func printSubnets(nis types.NetworkInterfaces, vpcs types.Vpcs, subnets types.Subnets, accounts types.Accounts) {
	table := NewTable()
	tableHeader := []string{"ACCOUNT ID", "AWS PROFILE", "VPC ID", "VPC NAME", "SUBNET ID", "SUBNET NAME", "CIDR", "OWNER ID", "OWNER PROFILE", "INTERFACES", "STATE"}
	table.AddRow(tableHeader...)
	for _, v := range subnets {
		var vpcName string
		if x := vpcs.GetById(v.VpcId); len(x) != 0 {
			vpcName = x[0].Name
		}
		numOfNis := len(nis.GetBySubnetId(v.SubnetId))

		table.AddRow(
			v.Account.Id,
			v.Account.Profile,
			v.VpcId,
			vpcName,
			v.SubnetId,
			v.Name,
			v.CidrBlock,
			v.OwnerId,
			accounts.GetById(v.OwnerId).Profile,
			out.FromInt(numOfNis),
			v.State,
		)
	}
	table.Print()
}

func findSubnets(arg string, subnets types.Subnets) types.Subnets {
	if IsIP(arg) {
		return subnets.GetByIp(arg)
	}
	if IsSubnetId(arg) {
		return subnets.GetById(arg)
	}
	if IsCIDR(arg) {
		return subnets.GetByCidr(arg)
	}
	fmt.Printf("argument %s can only be IP, CIDR, or network interface id", arg)
	os.Exit(1)
	return nil
}
