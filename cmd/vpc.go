package cmd

import (
	"fmt"
	"github.com/pete911/awf/internal/out"
	"github.com/pete911/awf/internal/types"
	"github.com/spf13/cobra"
	"os"
)

var (
	vpcCmd = &cobra.Command{
		Use:   "vpc",
		Short: "find vpc by IP or CIDR",
		Long:  "",
		Run:   runVpc,
	}
)

func init() {
	Root.AddCommand(vpcCmd)
}

func runVpc(cmd *cobra.Command, args []string) {
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

	var matched types.Vpcs
	for _, arg := range args {
		if found := findVpcs(arg, vpcs); len(found) > 0 {
			matched = append(matched, found...)
		}
	}

	if len(matched) == 0 {
		fmt.Printf("searched %d vpcs, but none matched\n", len(vpcs))
		return
	}

	printVpcs(matched, accounts)
}

func printVpcs(in types.Vpcs, accounts types.Accounts) {
	table := out.NewTable(os.Stdout)
	tableHeader := []string{"ACCOUNT ID", "AWS PROFILE", "ID", "NAME", "CIDR", "OWNER ID", "OWNER PROFILE", "STATE", "DEFAULT"}
	table.AddRow(tableHeader...)
	for _, v := range in {
		table.AddRow(
			v.Account.Id,
			v.Account.Profile,
			v.VpcId,
			out.TrimTo(v.Name, 40),
			v.CidrBlock,
			v.OwnerId,
			accounts.GetById(v.OwnerId).Profile,
			v.State,
			fmt.Sprintf("%t", v.IsDefault),
		)
	}
	table.Print()
}

func findVpcs(arg string, vpcs types.Vpcs) types.Vpcs {
	if IsIP(arg) {
		return vpcs.GetByIp(arg)
	}
	if IsEniId(arg) {
		return vpcs.GetById(arg)
	}
	if IsCIDR(arg) {
		return vpcs.GetByCidr(arg)
	}
	fmt.Printf("argument %s can only be IP, CIDR, or network interface id", arg)
	os.Exit(1)
	return nil
}
