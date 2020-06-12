module github.com/deta/deta-cli/api

go 1.13

replace github.com/deta/deta-cli/auth => ../auth

require (
	github.com/aws/aws-sdk-go v1.31.14
	github.com/deta/deta-cli/auth v0.0.0-00010101000000-000000000000
)
