package types

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"net/netip"
)

type Subnets []Subnet

func (v Subnets) GetByAccountId(in string) Subnets {
	var out Subnets
	for _, subnet := range v {
		if subnet.Account.Id == in {
			out = append(out, subnet)
		}
	}
	return out
}

func (v Subnets) GetById(in string) Subnets {
	var out Subnets
	for _, subnet := range v {
		if subnet.SubnetId == in {
			out = append(out, subnet)
		}
	}
	return out
}

func (v Subnets) GetByCidr(in string) Subnets {
	// we want to match 10.0.10.0/24 with 10.0.0.0/16 as well
	network, err := netip.ParsePrefix(in)
	if err != nil {
		return nil
	}

	var out Subnets
	for _, subnet := range v {
		subnetNetwork, err := netip.ParsePrefix(subnet.CidrBlock)
		if err != nil {
			return nil
		}
		if subnetNetwork.Overlaps(network) {
			out = append(out, subnet)
		}
	}
	return out
}

func (v Subnets) GetByIp(in string) Subnets {
	ip, err := netip.ParseAddr(in)
	if err != nil {
		return nil
	}

	var out Subnets
	for _, subnet := range v {
		network, err := netip.ParsePrefix(subnet.CidrBlock)
		if err != nil {
			out = append(out, subnet)
			continue
		}
		if network.Contains(ip) {
			out = append(out, subnet)
			continue
		}
	}
	return out
}

type Subnet struct {
	Account          Account
	Region           string
	SubnetId         string
	Name             string
	VpcId            string
	CidrBlock        string
	AvailabilityZone string
	OwnerId          string
	State            string
}

func ToSubnets(account Account, region string, in []types.Subnet) Subnets {
	var out Subnets
	for _, v := range in {
		out = append(out, ToSubnet(account, region, v))
	}
	return out
}

func ToSubnet(account Account, region string, in types.Subnet) Subnet {
	return Subnet{
		Account:          account,
		Region:           region,
		SubnetId:         aws.ToString(in.SubnetId),
		Name:             toTags(in.Tags)["Name"],
		VpcId:            aws.ToString(in.VpcId),
		CidrBlock:        aws.ToString(in.CidrBlock),
		AvailabilityZone: aws.ToString(in.AvailabilityZone),
		OwnerId:          aws.ToString(in.OwnerId),
		State:            string(in.State),
	}
}
