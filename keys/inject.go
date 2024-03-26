package keys

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2instanceconnect"
)

func AddMyKeyToEc2(ctx context.Context, instanceID, publicKey, region, profile, osUser, az string) error {
	// Load AWS configuration from environment variables, shared config, and shared credentials.
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithSharedConfigProfile(profile),
	)
	if err != nil {
		return err
	}

	// Create EC2 and EC2 Instance Connect clients
	//ec2Client := ec2.NewFromConfig(cfg)
	ec2InstanceConnectClient := ec2instanceconnect.NewFromConfig(cfg)

	// Prepare the request input
	input := &ec2instanceconnect.SendSSHPublicKeyInput{
		AvailabilityZone: &az,
		InstanceId:       &instanceID,
		InstanceOSUser:   &osUser,
		SSHPublicKey:     &publicKey,
	}

	// Send the SSH public key to the instance
	_, err = ec2InstanceConnectClient.SendSSHPublicKey(ctx, input)
	if err != nil {
		return err
	}

	//fmt.Println("SSH public key sent successfully", res.Success, " - ")
	return nil
}
