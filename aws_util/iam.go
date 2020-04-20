package aws_util

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/falmotlag/aws-bulk-cli/errors"
)

func AssumeIAmRole(roleArn string) (map[string]string, error) {
	sess, err := session.NewSession(&aws.Config{})
	if err != nil {
		return nil, errors.WithStackTrace(err)
	}

	_, err = sess.Config.Credentials.Get()
	if err != nil {
		return nil, errors.WithStackTraceAndPrefix(err, "Error finding AWS credentials")
	}

	stsClient := sts.New(sess)

	input := sts.AssumeRoleInput{
		RoleArn:         aws.String(roleArn),
		RoleSessionName: aws.String(fmt.Sprintf("aws-bulk-cli-%d", time.Now().UTC().UnixNano())),
	}

	StdOut, err := stsClient.AssumeRole(&input)
	if err != nil {
		return nil, errors.WithStackTrace(err)
	}

	creds := make(map[string]string)
	creds["AWS_ACCESS_KEY_ID"] = aws.StringValue(StdOut.Credentials.AccessKeyId)
	creds["AWS_SECRET_ACCESS_KEY"] = aws.StringValue(StdOut.Credentials.SecretAccessKey)
	creds["AWS_SESSION_TOKEN"] = aws.StringValue(StdOut.Credentials.SessionToken)
	creds["AWS_SECURITY_TOKEN"] = aws.StringValue(StdOut.Credentials.SessionToken)

	return creds, nil
}
