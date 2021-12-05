package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/eflem00/go-example-app/usecases"
	"github.com/rs/zerolog/log"
)

type QueueController struct {
	resultUsecase *usecases.ResultUsecase
}

func NewQueueController(resultUsecase *usecases.ResultUsecase) *QueueController {
	return &QueueController{
		resultUsecase,
	}
}

type MsgBody struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (controller *QueueController) process(msg *sqs.Message) error {
	body := MsgBody{}
	err := json.Unmarshal([]byte(*msg.Body), &body)

	if err != nil {
		return err
	}

	ctx := context.TODO() // TODO: Derive better context per msg

	log.Debug().Msg(fmt.Sprintf("%+v", body))

	return controller.resultUsecase.WriteResult(ctx, body.Key, body.Value)
}

func (controller *QueueController) Start() error {
	log.Info().Msg("Starting queue controller")

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	queueName := os.Getenv("QUEUE_NAME")

	result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queueName,
	})

	if err != nil {
		log.Err(err).Msg(fmt.Sprintf("Error fetching queue url for name %v", queueName))
		return err
	}

	log.Info().Msg(fmt.Sprintf("listening to %v", *result.QueueUrl))

	// poll for messages indefinitely
	for {
		msgResult, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
			AttributeNames: []*string{
				aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
			},
			MessageAttributeNames: []*string{
				aws.String(sqs.QueueAttributeNameAll),
			},
			QueueUrl:            result.QueueUrl,
			MaxNumberOfMessages: aws.Int64(10),
			VisibilityTimeout:   aws.Int64(30),
		})

		if err != nil {
			log.Err(err).Msg("error receiving message")
			continue
		}

		if len(msgResult.Messages) == 0 {
			continue
		}

		// normal sqs queues do not guarantee ordering so we should process multiple messages concurrently for better performance
		var wg sync.WaitGroup
		for _, message := range msgResult.Messages {
			wg.Add(1)

			go func(msg *sqs.Message) {
				defer wg.Done()

				log.Debug().Msg(fmt.Sprintf("processing message: %v %v", *msg.MessageId, *msg.Body))

				err := controller.process(msg)

				// if we fail to process the message we should not delete.
				// after a certain amount of reads without a delete the message will be automatically moved to a dead-letter queue
				if err != nil {
					log.Err(err).Msg(fmt.Sprintf("error processing message: %v", *msg.MessageId))
					return
				}

				_, err = svc.DeleteMessage(&sqs.DeleteMessageInput{
					QueueUrl:      result.QueueUrl,
					ReceiptHandle: msg.ReceiptHandle,
				})

				if err != nil {
					log.Err(err).Msg(fmt.Sprintf("error deleting message: %v", *msg.MessageId))
					return
				}

				log.Debug().Msg(fmt.Sprintf("processed message: %v", *msg.MessageId))
			}(message)

		}

		wg.Wait()
	}
}

func (controller *QueueController) Exit() {
	log.Error().Msg("detected exit in queue controller")
}
