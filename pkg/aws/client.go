package aws

import (
	"context"
	"log"
	"fmt"
	"errors"
)
import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type Client struct {
	config *aws.Config
}

func New() *Client {
	return &Client{config:nil}
}

func (c *Client) GetConfig() *aws.Config {
	return c.config
}

func (c *Client) SetConfig(region *string, account_id *string, role *string) (*aws.Config, error) {

	log.Println(c)

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
	sessionName := *role+"-progenitor"

	// check for cached config
	if c.config != nil {
		return c.config, nil
	}

	// new config	
	config, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}

	// Create the credentials from AssumeRoleProvider to assume the role
	// referenced by the "myRoleARN" ARN.
	stsSvc := sts.New(config)
	input := &sts.AssumeRoleInput{RoleArn: aws.String(arn), RoleSessionName: aws.String(sessionName)}
	out, err := stsSvc.AssumeRoleRequest(input).Send(context.TODO())
	if err != nil {
		return nil, err
	}

	awsConfig := stsSvc.Config.Copy()
	awsConfig.Credentials = CredentialsProvider{Credentials: out.Credentials}

	c.config = &awsConfig


	log.Println(c)

	return &awsConfig, nil
}

type CredentialsProvider struct {
	*sts.Credentials
}

func (s CredentialsProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	if s.Credentials == nil {
		return aws.Credentials{}, errors.New("sts credentials are nil")
	}

	return aws.Credentials{
		AccessKeyID:     aws.StringValue(s.AccessKeyId),
		SecretAccessKey: aws.StringValue(s.SecretAccessKey),
		SessionToken:    aws.StringValue(s.SessionToken),
		Expires:         aws.TimeValue(s.Expiration),
	}, nil
}
