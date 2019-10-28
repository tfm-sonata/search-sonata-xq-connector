package main

import (
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/search-sonata-xq-connector/awshelper"
	fraud "git-codecommit.eu-central-1.amazonaws.com/search-sonata-xq-connector/pStore"
	"git-codecommit.eu-central-1.amazonaws.com/search-sonata-xq-connector/search"
	"git-codecommit.eu-central-1.amazonaws.com/search-sonata-xq-connector/util"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"os"
	"time"
)

var (
	initialized = false
)

// Handler is our lambda handler invoked by the `lambda.Start` function call
//func Handler(ctx context.Context, event events.SQSEvent) error {
func Handler(payload map[string]interface{}) (interface{}, error) {

	// load all configs
	log.Println("is initialized : ", initialized)
	if !initialized {
		err := fraud.LoadServiceConfigFromParameterStore(os.Getenv("CONFIG_PROVIDER"), os.Getenv("SERVICENAME"))
		if err != nil {
			log.Println("Error reading config from parameter store:", err)
			return payload, err
		}
		log.Println("parameter store loaded ")
		awshelper.CreateSession()
		initialized = true
	}
	log.Println("initialized parameter store: ")

	// load all configs
	incomingJson := util.MapToString(payload)

	// if the request is a warmup request we stop here ...
	if pl := payload["source"]; pl != nil && pl.(string) == "serverless-warmup-request" {
		time.Sleep(95 * time.Millisecond)
		return os.Getenv("AWS_LAMBDA_LOG_STREAM_NAME"), nil
	}

	// if the request is a message, we need to handle inputString for handleMessage different
	if pl := payload["Records"]; pl != nil {
		//records := pl.(map[string]interface{})
		//messageBody := records["body"]
		//incomingJson = fmt.Sprintf(fmt.Sprintf("%v", messageBody))
		records := pl.([]interface{})
		messageBody := records[0].(map[string]interface{})["body"]
		incomingJson = fmt.Sprintf(fmt.Sprintf("%v", messageBody))
	}

	log.Println("Message received:", incomingJson)

	config := fraud.ParameterStoreConfig{
		ConfigProvider: os.Getenv("CONFIG_PROVIDER"),
		Servicename:    os.Getenv("SERVICENAME"),
	}
	//timeStart := time.Now()

	// defer awshelper.AddMetric(awshelper.RECEIVED_MESSAGES, awshelper.UNIT_COUNT, 1)
	tfmResponse, err := search.HandleMessage(config, incomingJson)
	if err != nil {
		log.Println("Error in handleMessage()", err)
		return payload, err
	}
	//timeEnd := time.Now()
	log.Println("SendMessage to Queue ", incomingJson)
	//timeEnd := time.Now()
	// defer awshelper.AddMetric(awshelper.RESPONSE_TIME, awshelper.UNIT_MILLISECONDS, float64(timeEnd.Sub(timeStart)/time.Millisecond))
	return tfmResponse, nil
}

func main() {
	lambda.Start(Handler)
}
