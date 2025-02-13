package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/pete911/awf/internal/types"
	"os"
	"sync"
	"time"
)

type importer func(account types.Account, cfg aws.Config, file File) error

var importers = []importer{
	ec2Import,
}

func Import(region string) error {
	cfg, err := newAwsConfig(region)
	if err != nil {
		return err
	}

	account, err := getCurrentAWSAccount(cfg)
	if err != nil {
		return err
	}
	file, err := initFile(account, cfg.Region)
	if err != nil {
		return err
	}

	var hasErrors bool
	var wg sync.WaitGroup
	for _, i := range importers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := i(account, cfg, file); err != nil {
				fmt.Println(err.Error())
				hasErrors = true
			}
		}()
	}
	wg.Wait()
	if hasErrors {
		return errors.New("import failed")
	}
	return nil
}

func newAwsConfig(awsRegion string) (aws.Config, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return aws.Config{}, fmt.Errorf("load aws config: %w", err)
	}

	if awsRegion == "" && cfg.Region == "" {
		return aws.Config{}, errors.New("missing aws region")
	}

	if awsRegion != "" {
		cfg.Region = awsRegion
	}
	return cfg, nil
}

func getCurrentAWSAccount(cfg aws.Config) (types.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stsSvc := sts.NewFromConfig(cfg)
	callerIdentity, err := stsSvc.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return types.Account{}, err
	}

	iamSvc := iam.NewFromConfig(cfg)
	accountAliases, err := iamSvc.ListAccountAliases(ctx, &iam.ListAccountAliasesInput{})
	if err != nil {
		return types.Account{}, err
	}

	var accountAlias string
	if len(accountAliases.AccountAliases) > 0 {
		accountAlias = accountAliases.AccountAliases[0]
	}

	return types.Account{
		Id:      aws.ToString(callerIdentity.Account),
		Profile: os.Getenv("AWS_PROFILE"),
		Alias:   accountAlias,
	}, nil
}
