package types

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"net/netip"
)

type Vpcs []Vpc

func (v Vpcs) GetByAccountId(in string) Vpcs {
	var out Vpcs
	for _, vpc := range v {
		if vpc.Account.Id == in {
			out = append(out, vpc)
		}
	}
	return out
}

func (v Vpcs) GetById(in string) Vpcs {
	var out Vpcs
	for _, vpc := range v {
		if vpc.VpcId == in {
			out = append(out, vpc)
		}
	}
	return out
}

func (v Vpcs) GetByCidr(in string) Vpcs {
	// we want to match 10.0.10.0/24 with 10.0.0.0/16 as well
	network, err := netip.ParsePrefix(in)
	if err != nil {
		return nil
	}

	var out Vpcs
	for _, vpc := range v {
		vpcNetwork, err := netip.ParsePrefix(vpc.CidrBlock)
		if err != nil {
			return nil
		}
		if vpcNetwork.Overlaps(network) {
			out = append(out, vpc)
		}
	}
	return out
}

func (v Vpcs) GetByIp(in string) Vpcs {
	ip, err := netip.ParseAddr(in)
	if err != nil {
		return nil
	}

	var out Vpcs
	for _, vpc := range v {
		network, err := netip.ParsePrefix(vpc.CidrBlock)
		if err != nil {
			out = append(out, vpc)
			continue
		}
		if network.Contains(ip) {
			out = append(out, vpc)
			continue
		}
	}
	return out
}

type Vpc struct {
	Account   Account
	Region    string
	VpcId     string
	Name      string
	CidrBlock string
	IsDefault bool
	OwnerId   string
	State     string
}

func ToVpcs(account Account, region string, in []types.Vpc) Vpcs {
	var out Vpcs
	for _, v := range in {
		out = append(out, ToVpc(account, region, v))
	}
	return out
}

func ToVpc(account Account, region string, in types.Vpc) Vpc {
	return Vpc{
		Account:   account,
		Region:    region,
		VpcId:     aws.ToString(in.VpcId),
		Name:      toTags(in.Tags)["Name"],
		CidrBlock: aws.ToString(in.CidrBlock),
		IsDefault: aws.ToBool(in.IsDefault),
		OwnerId:   aws.ToString(in.OwnerId),
		State:     string(in.State),
	}
}
