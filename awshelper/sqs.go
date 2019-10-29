package awshelper

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"git-codecommit.eu-central-1.amazonaws.com/search-sonata-xq-connector/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
)

type TrackingMessage struct {
	Timestamp   int64  `json:"timestamp"`
	ServiceName string `json:"serviceName"`
	Action      string `json:"action"`
	SessionId   string `json:"sessionId"`
	RequestId   string `json:"requestId"`
	ContentType string `json:"contentType"`
	BodyZ       string `json:"bodyZ"`
	Title       string `json:"title"`
	ActionType  string `json:"actionType"`
}

var (
	sqsClient *sqs.SQS
)

func SendMessage(queueUrl, messageBody string) error {

	if sqsClient == nil {
		sqsClient = sqs.New(AWSSession())
	}

	log.Println("Sending message:", messageBody)
	_, err := sqsClient.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(0),
		MessageBody:  aws.String(messageBody),
		QueueUrl:     &queueUrl,
	})

	if err != nil {
		log.Println("Couldn't send message via SQS:", messageBody, err)
		return err
	}

	return nil
}

func SendTrackingMessage(trackingMessage TrackingMessage) error {

	bodyZ, err := compressMessage(trackingMessage.BodyZ)
	if err != nil {
		log.Println("Couldn't compress messageBody of trackingMessage", err)
		return err
	}
	trackingMessage.BodyZ = bodyZ

	byteJson, err := json.Marshal(trackingMessage)
	if err != nil {
		log.Println("Couldn't marshal trackingMessage into json", err)
		return err
	}
	json := string(byteJson)

	err = SendMessage(model.RequestConfig.SqsBaseUrl+model.RequestConfig.QueueNameTracking, json)
	if err != nil {
		log.Println("Couldn't send message", err)
		return err
	}

	return nil
}

// compress the stringified json and create a base64 zLib string
func compressMessage(msg string) (string, error) {
	var buf bytes.Buffer
	gz := zlib.NewWriter(&buf)
	if _, err := gz.Write([]byte(msg)); err != nil {
		return "", err
	}
	if err := gz.Flush(); err != nil {
		return "", err
	}
	if err := gz.Close(); err != nil {
		return "", err
	}
	// base64 encoding for the zippedBytes.Buffer
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())

	return encoded, nil
}
