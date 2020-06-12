package api

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

var (
	accessKey string
	secretKey string
)

// DeployLambda xx
func (c *DetaClient) DeployLambda(programID string, zipFile []byte) error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("eu-central-1"),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
	if err != nil {
		return err
	}
	fx := lambda.New(sess)
	_, err = fx.UpdateFunctionCode(&lambda.UpdateFunctionCodeInput{
		Publish:      aws.Bool(true),
		FunctionName: aws.String(programID),
		ZipFile:      zipFile,
	})
	if err != nil {
		return err
	}
	return nil
}
