package consumer

import (
  "fmt"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/credentials"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/sqs"
)


func NewSQSClient(key string, secret string, region string) (*sqs.SQS, error) {

  awsConfig := &aws.Config {
    Credentials: credentials.NewStaticCredentials(key, secret, ""),
    Region: aws.String(region),
  }

  sess, err := session.NewSession()
  if err != nil {
    fmt.Println("failed to create session,", err)
    return nil, err
  }
  svc := sqs.New(sess, awsConfig)

  return svc, nil
}