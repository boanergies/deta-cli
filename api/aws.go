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

// LambdaClient xx
type LambdaClient struct {
	fx *lambda.Lambda
}

// NewLambdaClient xx
func NewLambdaClient() *LambdaClient {
	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String("eu-central-1"),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
	fx := lambda.New(sess)
	return &LambdaClient{
		fx: fx,
	}
}

// DeployLambda xx
func (c *LambdaClient) DeployLambda(programID string, zipFile []byte) error {
	_, err := c.fx.UpdateFunctionCode(&lambda.UpdateFunctionCodeInput{
		//Publish:      aws.Bool(true),
		FunctionName: aws.String(programID),
		ZipFile:      zipFile,
	})
	if err != nil {
		return err
	}
	return nil
}
