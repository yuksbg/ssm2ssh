package profiles

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
)

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
