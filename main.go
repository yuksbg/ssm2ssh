package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/mmmorris1975/ssm-session-client/ssmclient"
	lop "github.com/samber/lo/parallel"
	"github.com/yuksbg/ssm2ssh/keys"
	"github.com/yuksbg/ssm2ssh/profiles"
	"log"
	"net"
	"os"
)

func main() {

	instanceID := os.Args[1]
	foundInstance := types.Instance{}
	foundProfile := ""
	useIAMProfile := false
	iamProfile := os.Getenv(`USE_IAM_ROLE`)

	if instanceID == "" {
		fmt.Println("All params are required")
		os.Exit(1)
	}
	if os.Getenv("USE_IAM_ROLE") != "" { /// try with aim
		useIAMProfile = true
		arnRole := os.Getenv("USE_IAM_ROLE")
		if arnRole == "1" {
			arnRole = ""
		}
		iamProfile = arnRole
		f, err := profiles.FindIAM(arnRole, instanceID)
		if err == nil {
			if *f.Placement.AvailabilityZone != "" {
				foundInstance = f
			}
		} else {
			return
		}
	} else {
		prf, _ := profiles.ListAWSProfiles()
		lop.ForEach(prf, func(item string, index int) {
			f, err := profiles.Find(item, instanceID)
			if err != nil {
				return
			}
			if *f.Placement.AvailabilityZone != "" {
				foundInstance = f
				foundProfile = item
			}
		})
	}

	if foundInstance.InstanceId == nil {
		os.Exit(1)
	}

	// get default key
	sshKey, err := keys.GetDefaultSSHKey()
	if err != nil {
		fmt.Println("Default SSH key not found - ", err.Error())
		os.Exit(1)
	}
	if err := keys.AddMyKeyToEc2(context.Background(),
		*foundInstance.InstanceId,
		sshKey,
		"eu-central-1",
		foundProfile,
		"root",
		*foundInstance.Placement.AvailabilityZone); err != nil {

		fmt.Println("SSH Add error - ", err.Error())
		os.Exit(1)
	}
	var cfg aws.Config
	if useIAMProfile {
		cfg = profiles.GetIAMConfig(iamProfile)
	} else {
		cfg, err = config.LoadDefaultConfig(context.Background(), config.WithSharedConfigProfile(foundProfile))
		if err != nil {
			log.Fatal(err)
		}
	}

	var port int
	t, p, err := net.SplitHostPort(instanceID)
	if err == nil {
		port, err = net.LookupPort("tcp", p)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		t = instanceID
	}

	tgt, err := ssmclient.ResolveTarget(t, cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(ssmclient.SSHPluginSession(cfg, &ssmclient.PortForwardingInput{Target: tgt, LocalPort: port}))
}
