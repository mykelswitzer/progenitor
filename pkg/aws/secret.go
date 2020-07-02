package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func (c *Client) GetSecret(secret string) (*GetSecretValueResponse, error) {

	svc := secretsmanager.New(c.GetConfig())	
	input := &secretsmanager.GetSecretValueInput{
	    SecretId: aws.String(secret),
	}

	req := svc.GetSecretValueRequest(input)
	result, err := req.Send(context.Background())
	if err != nil {
	    if aerr, ok := err.(awserr.Error); ok {
	        switch aerr.Code() {
	        case secretsmanager.ErrCodeResourceNotFoundException:
	            fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
	        case secretsmanager.ErrCodeInvalidParameterException:
	            fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
	        case secretsmanager.ErrCodeInvalidRequestException:
	            fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
	        case secretsmanager.ErrCodeDecryptionFailure:
	            fmt.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
	        case secretsmanager.ErrCodeInternalServiceError:
	            fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
	        default:
	            fmt.Println(aerr.Error())
	        }
	        return nil, aerr
	    } else {
	        fmt.Println(err.Error())
	        return nil, err
	    }
	}

	return result, nil

}