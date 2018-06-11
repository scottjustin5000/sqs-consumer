# SQS-Consumer

Simple long polling sqs-consumer.

### Usage

```go
import (
  "github.com/scottjustin5000/sqs-consumer/consumer"
)


func main() {

  svc, err := consumer.NewSQSClient(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), os.Getenv("AWS_REGION"))

  if err != nil {
    fmt.Println(err.Error())
    return
  }

  er := consumer.Init(svc, "some-sqs-queue-stage")

  if er != nil {
    fmt.Println(er.Error())
    return
  }


  consumer.Start(consumer.CallbackFunc(func(msg *sqs.Message) error {
    fmt.Println(aws.StringValue(msg.Body))
    return nil
  }))
}

```
