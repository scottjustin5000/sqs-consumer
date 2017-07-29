package consumer

import (
  "fmt"
  "log"
  "sync"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/sqs"
)

var queueUrl string = ""
var svc *sqs.SQS

//public vars

//number of messages to retreive, 1-10
var NumberOfMessages int64 = 10
// duration to wait for a message to show up
var WaitTime int64 = 20


type CallbackFunc func(msg *sqs.Message) error

func (f CallbackFunc) MessageHandler(msg *sqs.Message) error {
  return f(msg)
}

type Handler interface {
  MessageHandler(msg *sqs.Message) error
}

func Init(sv *sqs.SQS, name string, ) error {
  svc = sv
  params := &sqs.GetQueueUrlInput{
    QueueName: aws.String(name), 
  }
  resp, err := svc.GetQueueUrl(params)
  queueUrl = aws.StringValue(resp.QueueUrl)

  return err
}

func Start(cb Handler) {

  for {
    fmt.Println("Listening...")
    params := &sqs.ReceiveMessageInput {
      QueueUrl: aws.String(queueUrl),
      MaxNumberOfMessages: aws.Int64(NumberOfMessages),
      MessageAttributeNames: []*string {
        aws.String("All"), 
      },
      WaitTimeSeconds: aws.Int64(WaitTime),
    }

    resp, err := svc.ReceiveMessage(params)
    if err != nil {
      fmt.Println(err)
      continue
    }
    if len(resp.Messages) > 0 {
      processBatch(cb, resp.Messages)
    }
  }
}


func processBatch(cb Handler, messages []*sqs.Message) {
  numMessages := len(messages)

  var wg sync.WaitGroup
  wg.Add(numMessages)
  for i := range messages {
    go func(m *sqs.Message) {
      defer wg.Done()
      if err := processMessage(m, cb); err != nil {
        log.Printf("message error -- %s", err.Error())
      }
    }(messages[i])
  }

  wg.Wait()
}

func processMessage(m *sqs.Message, cb Handler) error {
  fmt.Println("processing message...")
  var err error
  err = cb.MessageHandler(m)
  
  if err != nil {
    return err
  }

  params := &sqs.DeleteMessageInput{
    QueueUrl:      aws.String(queueUrl), 
    ReceiptHandle: m.ReceiptHandle,      
  }
  _, err = svc.DeleteMessage(params)
  if err != nil {
    return err
  }
  fmt.Println("removing message: %s", aws.StringValue(m.ReceiptHandle))

  return nil
}