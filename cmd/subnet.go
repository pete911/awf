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

	printSubnets(matched, vpcs, accounts)
}

func printSubnets(in types.Subnets, vpcs types.Vpcs, accounts types.Accounts) {
	table := out.NewTable(os.Stdout)
	tableHeader := []string{"ACCOUNT ID", "AWS PROFILE", "VPC ID", "VPC NAME", "SUBNET ID", "SUBNET NAME", "CIDR", "OWNER ID", "OWNER PROFILE", "STATE"}
	table.AddRow(tableHeader...)
	for _, v := range in {
		var vpcName string
		if x := vpcs.GetById(v.VpcId); len(x) != 0 {
			vpcName = x[0].Name
		}

		table.AddRow(
			v.Account.Id,
			v.Account.Profile,
			v.VpcId,
			out.TrimTo(vpcName, 40),
			v.SubnetId,
			out.TrimTo(v.Name, 40),
			v.CidrBlock,
			v.OwnerId,
			accounts.GetById(v.OwnerId).Profile,
			v.State,
		)
	}
	table.Print()
}

func findSubnets(arg string, subnets types.Subnets) types.Subnets {
	if IsIP(arg) {
		return subnets.GetByIp(arg)
	}
	if IsEniId(arg) {
		return subnets.GetById(arg)
	}
	if IsCIDR(arg) {
		return subnets.GetByCidr(arg)
	}
	fmt.Printf("argument %s can only be IP, CIDR, or network interface id", arg)
	os.Exit(1)
	return nil
}
