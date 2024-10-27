package store

import (
	"encoding/json"
	"errors"
	"fmt"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/pete911/awf/internal/types"
	"os"
	"path/filepath"
)

const (
	rootDir     = ".awf"
	accountFile = "_account"

	ec2VpcsKey              = "ec2.describe-vpcs"
	ec2SubnetsKey           = "ec2.describe-subnets"
	ec2NetworkInterfacesKey = "ec2.describe-network-interfaces"
)

type File struct {
	dir string
}

func LoadFile() (File, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return File{}, err
	}

	return File{
		dir: filepath.Join(home, rootDir),
	}, nil
}

func initFile(account types.Account, region string) (File, error) {
	f, err := LoadFile()
	if err != nil {
		return File{}, err
	}

	if err := os.MkdirAll(filepath.Join(f.dir, account.Id, region), 0755); err != nil {
		return File{}, err
	}

	if err := f.write(account, "", accountFile, account); err != nil {
		return File{}, fmt.Errorf("write account data: %w", err)
	}
	return f, nil
}

func (f File) ListAccounts() (types.Accounts, error) {
	entry, err := os.ReadDir(f.dir)
	if err != nil {
		return nil, err
	}

	var accounts []types.Account
	for _, e := range entry {
		if e.IsDir() {
			var account types.Account
			if err := f.read(filepath.Join(f.dir, e.Name(), accountFile), &account); err != nil {
				return nil, err
			}
			accounts = append(accounts, account)
		}
	}
	return accounts, nil
}

func (f File) ListRegions(account types.Account) ([]string, error) {
	entry, err := os.ReadDir(filepath.Join(f.dir, account.Id))
	if err != nil {
		return nil, err
	}

	var regions []string
	for _, e := range entry {
		if e.IsDir() {
			regions = append(regions, e.Name())
		}
	}
	return regions, nil
}

// DescribeNetworkInterfaces returns network interfaces. If the file is not found,
// NotFoundError is returned. Meaning that the user need to run import first.
func (f File) DescribeNetworkInterfaces() (types.NetworkInterfaces, error) {
	accounts, err := f.ListAccounts()
	if err != nil {
		return nil, err
	}

	var networkInterfaces types.NetworkInterfaces
	for _, account := range accounts {
		regions, err := f.ListRegions(account)
		if err != nil {
			return nil, err
		}
		for _, region := range regions {
			path := f.filePath(account.Id, region, ec2NetworkInterfacesKey)

			var nis []ec2types.NetworkInterface
			if err := f.read(path, &nis); err != nil {
				return nil, err
			}
			networkInterfaces = append(networkInterfaces, types.ToNetworkInterfaces(account, region, nis)...)
		}
	}
	return networkInterfaces, nil
}

// DescribeVpcs returns VPCs. If the file is not found,
// NotFoundError is returned. Meaning that the user need to run import first.
func (f File) DescribeVpcs() (types.Vpcs, error) {
	accounts, err := f.ListAccounts()
	if err != nil {
		return nil, err
	}

	var vpcs types.Vpcs
	for _, account := range accounts {
		regions, err := f.ListRegions(account)
		if err != nil {
			return nil, err
		}
		for _, region := range regions {
			path := f.filePath(account.Id, region, ec2VpcsKey)

			var awsVpcs []ec2types.Vpc
			if err := f.read(path, &awsVpcs); err != nil {
				return nil, err
			}
			vpcs = append(vpcs, types.ToVpcs(account, region, awsVpcs)...)
		}
	}
	return vpcs, nil
}

// DescribeSubnets returns subnets. If the file is not found,
// NotFoundError is returned. Meaning that the user need to run import first.
func (f File) DescribeSubnets() (types.Subnets, error) {
	accounts, err := f.ListAccounts()
	if err != nil {
		return nil, err
	}

	var subnets types.Subnets
	for _, account := range accounts {
		regions, err := f.ListRegions(account)
		if err != nil {
			return nil, err
		}
		for _, region := range regions {
			path := f.filePath(account.Id, region, ec2SubnetsKey)

			var awsSubnets []ec2types.Subnet
			if err := f.read(path, &awsSubnets); err != nil {
				return nil, err
			}
			subnets = append(subnets, types.ToSubnets(account, region, awsSubnets)...)
		}
	}
	return subnets, nil
}

func (f File) read(path string, v any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return NewNotFoundError(fmt.Sprintf("read: %s file does not exist, empty content", path))
		}
		return err
	}

	if err := json.Unmarshal(b, &v); err != nil {
		return fmt.Errorf("unmarshal %s: %w", path, err)
	}
	return nil
}

// Write writes content of the supplied (json) struct under supplied <name> file. Region
// can be empty (e.g. route53)
func (f File) write(account types.Account, region, name string, v any) error {
	path := f.filePath(account.Id, region, name)
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal %s: %w", path, err)
	}
	return os.WriteFile(path, b, 0644)
}

func (f File) filePath(accountId, region, name string) string {
	if region == "" {
		return filepath.Join(f.dir, accountId, name)
	}
	return filepath.Join(f.dir, accountId, region, name)
}
