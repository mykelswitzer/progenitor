package aws

import (
	"log"
	"fmt"
	"errors"
)
import (
	"github.com/aws/aws-sdk-go/"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
)

type Client struct {
	config *aws.Config
}

func (c Client) GetConfig() *aws.Config {
	return c.config
}

func (c Client) SetConfig(region *string, account_id *string, role *string) (*aws.Config, error) {

	// return no config for nil inputs
	if account_id == nil || region == nil || role == nil {
		err := errors.New("Region, Account Id, and Role Name to assume are all required.");
		log.Println(err.Error())
		return nil, err
	}

	arn := fmt.Sprintf(
		"arn:aws:iam::%v:role/%v",
		*account_id,
		*role,
	)

	// include region in cache key otherwise concurrency errors
	key := fmt.Sprintf("%v::%v", *region, arn)

	// check for cached config
	if c.configs != nil && c.configs[key] != nil {
		return c.configs[key]
	}

	// new config
	config, err := external.LoadDefaultAWSConfig()

	// Create the credentials from AssumeRoleProvider to assume the role
	// referenced by the "myRoleARN" ARN.
	stsSvc := sts.New(cfg)
	stsCredProvider := stscreds.NewAssumeRoleProvider(stsSvc, arn)

	config.Credentials = aws.NewCredentials(stsCredProvider)

	c.config = config

	return config, nil
}
