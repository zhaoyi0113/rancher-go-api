package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
)

var region = "ap-southeast-2"

var entityUuid = "c1b5baaf-0301-4bc6-9ce8-ae9cf919638f"
var siteUuid = "3c4ed0c6-a028-4c9f-a830-84c7cdde5db9"
var deviceUuid = "a37317aa-e17c-4831-a503-9f40546c13f6"

type TransactionRequest struct {
	Amount int
}

type Amount struct {
	Currency string `json:"currency"`
	Value    string `json:"value"`
}

type TransactionInitiated struct {
	Amount          Amount `json:"amount"`
	TransactionUuid string `json:"transactionUuid"`
	Timestamp       string `json:"timestamp"`
	TimestampUtc    string `json:"timestampUtc"`
	EntityUuid      string `json:"entityUuid"`
	SiteUuid        string `json:"siteUuid"`
	DeviceUuid      string `json:"deviceUuid"`
	Type            string `json:"type"`
	Scheme          string `json:"scheme"`
	SurchargeAmount Amount `json:"surchargeAmount"`
	SaleAmount      Amount `json:"saleAmount"`
	TipAmount       Amount `json:"tipAmount"`
}

type TransactionApproved struct {
	TransactionUuid     string `json:"transactionUuid"`
	EntityUuid          string `json:"entityUuid"`
	SiteUuid            string `json:"siteUuid"`
	DeviceUuid          string `json:"deviceUuid"`
	Timestamp           string `json:"timestamp"`
	TimestampUtc        string `json:"timestampUtc"`
	ResponseCode        string `json:"responseCode"`
	ResponseDescription string `json:"responseDescription"`
	ApprovalCode        string `json:"approvalCode"`
	Rrn                 string `json:"rrn"`
}

type DomainEvent struct {
	AggregateId string                 `json:"aggregateId"`
	EventId     string                 `json:"eventId"`
	Uri         string                 `json:"uri"`
	Aggregate   string                 `json:"aggregate"`
	Payload     map[string]interface{} `json:"payload"`
}

type AWSCredential struct {
	KeyId  string `json:"AWS_ACCESS_KEY_ID"`
	Secret string `json:"AWS_SECRET_ACCESS_KEY"`
}

func ProcessTransactionRequest(request []byte) {
	var transactionRequest TransactionRequest
	err := json.Unmarshal(request, &transactionRequest)
	if err != nil {
		log.Println("Failed to process transaction request", err)
		return
	}
	awsCredentials := getCredentials()
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsCredentials.KeyId, awsCredentials.Secret, "")),
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

func getTransactionScheme() string {
	awsKey := os.Getenv("CLOUD_PROVIDER")
	if awsKey == "AWS" {
		return "VISA"
	}
	return "MC"
}

func getInitiatedEvent(id string, amount int) string {
	initiated := TransactionInitiated{
		TransactionUuid: id,
		Amount: Amount{
			Currency: "AUD", Value: strconv.Itoa(amount),
		},
		SurchargeAmount: Amount{
			Currency: "AUD", Value: "0",
		},
		TipAmount: Amount{
			Currency: "AUD", Value: "0",
		},
		Timestamp:    time.Now().Format(time.RFC3339),
		TimestampUtc: time.Now().UTC().Format("2006-01-02T15:04:05Z07:00"),
		Type:         "PURCHASE",
		EntityUuid:   entityUuid,
		SiteUuid:     siteUuid,
		DeviceUuid:   deviceUuid,
		Scheme:       getTransactionScheme(),
	}
	bytes, _ := json.Marshal(initiated)
	payload := make(map[string]interface{})
	json.Unmarshal(bytes, &payload)
	event := DomainEvent{
		EventId:     uuid.NewString(),
		AggregateId: id,
		Aggregate:   "Transaction",
		Uri:         "pgs.Transaction.Initiated",
		Payload:     payload,
	}
	bytes, _ = json.Marshal(event)
	fmt.Println("send initiated transaction to eventbus", initiated.TransactionUuid)
	return string(bytes)
}

func getApprovedEvent(id string) string {
	approved := TransactionApproved{
		TransactionUuid:     id,
		EntityUuid:          entityUuid,
		SiteUuid:            siteUuid,
		DeviceUuid:          deviceUuid,
		Timestamp:           time.Now().Format(time.RFC3339),
		TimestampUtc:        time.Now().UTC().Format("2006-01-02T15:04:05Z07:00"),
		ApprovalCode:        "00",
		ResponseCode:        "00",
		ResponseDescription: "Approved",
	}
	bytes, _ := json.Marshal(approved)
	payload := make(map[string]interface{})
	json.Unmarshal(bytes, &payload)
	event := DomainEvent{
		EventId:     uuid.NewString(),
		AggregateId: id,
		Aggregate:   "Transaction",
		Uri:         "pgs.Transaction.Approved",
		Payload:     payload,
	}
	bytes, _ = json.Marshal(event)
	fmt.Println("send initiated transaction to eventbus", approved.TransactionUuid)
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
	fmt.Println("send event response", output.FailedEntryCount, output.Entries[0].ErrorCode, output.Entries[0].ErrorMessage)
}

func getCredentials() AWSCredential {
	var credentials AWSCredential
	credentials.KeyId = os.Getenv("AWS_ACCESS_KEY_ID")
	credentials.Secret = os.Getenv("AWS_SECRET_ACCESS_KEY")
	return credentials
}
