package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
)

var entityUuid = "5262a830-7b7a-41d4-ba81-062459d36992"
var siteUuid = "3c4ed0c6-a028-4c9f-a830-84c7cdde5db9"

type TransactionRequest struct {
	Amount int
}

type TransactionInitiated struct {
	Amount          int    `json:"amount"`
	TransactionUuid string `json:"transactionUuid"`
	Timestamp       string `json:"timestamp"`
	TimestampUtc    string `json:"timestampUtc"`
	EntityUuid      string `json:"entityUuid"`
	SiteUuid        string `json:"siteUuid"`
	Type            string `json:"type"`
}

type TransactionApproved struct {
	TransactionUuid     string `json:"transactionUuid"`
	EntityUuid          string `json:"entityUuid"`
	SiteUuid            string `json:"siteUuid"`
	Timestamp           string `json:"timestamp"`
	TimestampUtc        string `json:"timestampUtc"`
	ResponseCode        string `json:"responseCode"`
	ResponseDescription string `json:"responseDescription"`
	ApprovalCode        string `json:"approvalCode"`
	Rrn                 string `json:"rrn"`
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

	transactionUuid := uuid.NewString()
	initiatedEvent := getInitiatedEvent(transactionUuid, transactionRequest.Amount)
	approvedEvent := getApprovedEvent(transactionUuid)

	sendEvent(client, "pgs.Transaction.Initiated", initiatedEvent)
	sendEvent(client, "pgs.Transaction.Approved", approvedEvent)
}

func getInitiatedEvent(id string, amount int) string {
	initiated := TransactionInitiated{
		TransactionUuid: id,
		Amount:          amount,
		Timestamp:       time.Now().Format(time.RFC3339),
		TimestampUtc:    time.Now().UTC().Format("2006-01-02T15:04:05Z07:00"),
		Type:            "PURCHASE",
		EntityUuid:      entityUuid,
		SiteUuid:        siteUuid,
	}
	bytes, _ := json.Marshal(initiated)
	fmt.Println("send initiated transaction to eventbus", initiated.TransactionUuid)
	return string(bytes)
}

func getApprovedEvent(id string) string {
	approved := TransactionApproved{
		TransactionUuid: id,
		EntityUuid:      entityUuid,
		SiteUuid:        siteUuid,
		Timestamp:       time.Now().Format(time.RFC3339),
		TimestampUtc:    time.Now().UTC().Format("2006-01-02T15:04:05Z07:00"),
	}
	fmt.Println("send approved transaction to eventbus", approved.TransactionUuid)
	bytes, _ := json.Marshal(approved)
	return string(bytes)
}

func sendEvent(client *eventbridge.Client, eventUri string, event string) {
	fmt.Println("send event to event bus", event)
	output, err := client.PutEvents(context.TODO(), &eventbridge.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{
			{
				EventBusName: aws.String("dev-eventBus-global"),
				Detail:       aws.String(event),
				DetailType:   aws.String(eventUri),
				Source:       aws.String("pgs"),
			},
		},
	})
	if err != nil {
		log.Println("Failed to send event to event bridge", err)
		return
	}
	fmt.Println("send event response", output.FailedEntryCount, output.Entries)
}
