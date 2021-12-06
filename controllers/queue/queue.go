package queue

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/eflem00/go-example-app/common"
	"github.com/eflem00/go-example-app/usecases"
)

type QueueController struct {
	settings      *common.Settings
	logger        *common.Logger
	resultUsecase *usecases.ResultUsecase
}

func NewQueueController(settings *common.Settings, logger *common.Logger, resultUsecase *usecases.ResultUsecase) *QueueController {
	return &QueueController{
		settings,
		logger,
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

	return controller.resultUsecase.WriteResult(ctx, body.Key, body.Value)
}

func (controller *QueueController) Start() error {
	controller.logger.Info("Starting queue controller")

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	queueName := controller.settings.QueueName

	result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queueName,
	})

	if err != nil {
		controller.logger.Errf(err, "Error fetching queue url for name %v", queueName)
		return err
	}

	controller.logger.Infof("listening to %v", *result.QueueUrl)

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
			controller.logger.Err(err, "error receiving message")
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

				controller.logger.Debugf("starting message: %v %v", *msg.MessageId, *msg.Body)

				err := controller.process(msg)

				// if we fail to process the message we should not delete.
				// after a certain amount of reads without a delete the message will be automatically moved to a dead-letter queue
				if err != nil {
					controller.logger.Errf(err, "error processing message: %v", *msg.MessageId)
					return
				}

				_, err = svc.DeleteMessage(&sqs.DeleteMessageInput{
					QueueUrl:      result.QueueUrl,
					ReceiptHandle: msg.ReceiptHandle,
				})

				if err != nil {
					controller.logger.Errf(err, "error deleting message: %v", *msg.MessageId)
					return
				}

				controller.logger.Debugf("finished message: %v", *msg.MessageId)
			}(message)

		}

		wg.Wait()
	}
}

func (controller *QueueController) Exit() {
	controller.logger.Error("detected exit in queue controller")
}
