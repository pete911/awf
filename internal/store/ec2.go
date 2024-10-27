package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/pete911/awf/internal/types"
	"sync"
	"time"
)

type ec2Importer func(svc *ec2.Client) (string, any, error)

var ec2Importers = []ec2Importer{
	describeVpcs,
	describeSubnets,
	describeNetworkInterfaces,
}

func ec2Import(account types.Account, cfg aws.Config, file File) error {
	svc := ec2.NewFromConfig(cfg)

	var hasErrors bool
	var wg sync.WaitGroup
	for _, i := range ec2Importers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			name, content, err := i(svc)
			if err != nil {
				fmt.Println(err.Error())
				hasErrors = true
				return
			}
			if err := file.write(account, cfg.Region, name, content); err != nil {
				fmt.Println(err.Error())
				hasErrors = true
			}
		}()
	}
	wg.Wait()
	if hasErrors {
		return errors.New("ec2 import failed")
	}
	return nil
}

func describeVpcs(svc *ec2.Client) (string, any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var vpcs []ec2types.Vpc
	in := &ec2.DescribeVpcsInput{}
	for {
		out, err := svc.DescribeVpcs(ctx, in)
		if err != nil {
			return "", nil, err
		}
		vpcs = append(vpcs, out.Vpcs...)
		if aws.ToString(out.NextToken) == "" {
			break
		}
		in.NextToken = out.NextToken
	}
	return ec2VpcsKey, vpcs, nil
}

func describeSubnets(svc *ec2.Client) (string, any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var subnets []ec2types.Subnet
	in := &ec2.DescribeSubnetsInput{}
	for {
		out, err := svc.DescribeSubnets(ctx, in)
		if err != nil {
			return "", nil, err
		}
		subnets = append(subnets, out.Subnets...)
		if aws.ToString(out.NextToken) == "" {
			break
		}
		in.NextToken = out.NextToken
	}
	return ec2SubnetsKey, subnets, nil
}

func describeNetworkInterfaces(svc *ec2.Client) (string, any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var nis []ec2types.NetworkInterface
	in := &ec2.DescribeNetworkInterfacesInput{}
	for {
		out, err := svc.DescribeNetworkInterfaces(ctx, in)
		if err != nil {
			return "", nil, err
		}
		nis = append(nis, out.NetworkInterfaces...)
		if aws.ToString(out.NextToken) == "" {
			break
		}
		in.NextToken = out.NextToken
	}
	return ec2NetworkInterfacesKey, nis, nil
}
