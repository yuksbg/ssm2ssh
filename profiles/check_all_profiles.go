package profiles

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
)

func GetIAMConfig(AssumeRoleARN string) aws.Config {
	defRegion := `eu-central-1`
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(defRegion), config.WithSharedConfigProfile(``))
	if err != nil {
		return aws.Config{}
	}
	if AssumeRoleARN != "" {
		ctx := context.TODO()
		roleArn := AssumeRoleARN
		awsRegion := defRegion
		cfg, err = assumeRole(ctx, roleArn, awsRegion)
	}
	return cfg
}

func FindIAM(AssumeRoleARN string, instanceName string) (resultInstance types.Instance, err error) {

	cfg := GetIAMConfig(AssumeRoleARN)

	// Create an EC2 client
	client := ec2.NewFromConfig(cfg)

	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceName},
	}

	resp, err := client.DescribeInstances(context.Background(), input)
	if err != nil {
		return
	}
	// Check if the instance is found
	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			resultInstance = instance
		}
	}
	return
}

func assumeRole(ctx context.Context, roleArn string, awsRegion string) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(awsRegion))
	if err != nil {
		return aws.Config{}, err
	}

	// Create a new STS client.
	stsClient := sts.NewFromConfig(cfg)

	// Construct the credentials provider
	creds := stscreds.NewAssumeRoleProvider(stsClient, roleArn)

	// Create a new AWS config with the assumed role's credentials
	assumedRoleCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(awsRegion), config.WithCredentialsProvider(creds))
	if err != nil {
		return aws.Config{}, err
	}

	return assumedRoleCfg, nil
}

func Find(profile string, instanceName string) (resultInstance types.Instance, err error) {
	// Load AWS configuration from the default credential chain
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithSharedConfigProfile(profile))
	if err != nil {
		return
	}

	// Create an EC2 client
	client := ec2.NewFromConfig(cfg)

	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceName},
	}

	resp, err := client.DescribeInstances(context.Background(), input)
	if err != nil {
		return
	}
	// Check if the instance is found
	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			resultInstance = instance
		}
	}
	return
}

func ListAWSProfiles() ([]string, error) {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return []string{}, err
	}

	sshDir := filepath.Join(homeDir, ".aws")
	publicKeyPath := filepath.Join(sshDir, "credentials")

	cfg, err := ini.Load(publicKeyPath)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
	}

	// Get a list of all sections
	sections := cfg.SectionStrings()

	// Print all section names
	ses := []string{}
	for _, section := range sections {
		ses = append(ses, section)
	}
	return ses, nil
}
