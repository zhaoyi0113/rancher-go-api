package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/aws/aws-sdk-go/aws"
)

type TransactionRequest struct {
	Amount int
}

func ProcessTransactionRequest(request []byte) {
	var transactionRequest TransactionRequest
	err := json.Unmarshal(request, &transactionRequest)
	if err != nil {
		log.Println("Failed to process transaction request", err)
		return
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("ap-southeast-2"),
	)
	if err != nil {
		log.Println("unable to load SDK config", err)
		return
	}
	client := eventbridge.NewFromConfig(cfg)

	output, err := client.PutEvents(context.TODO(), &eventbridge.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{
			{
				EventBusName: aws.String("dev-eventBus-global"),
				Detail:       aws.String("{ \"key1\": \"value1\", \"key2\": \"value2\" }"),
				DetailType:   aws.String("pgs.Transaction.Initiated"),
				Source:       aws.String("pgs"),
			},
		},
	})
	if err != nil {
		log.Println("Failed to send event to event bridge", err)
		return
	}
	fmt.Println("send event response", output)
}
