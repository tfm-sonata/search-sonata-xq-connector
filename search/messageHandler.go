package search

import (
	"encoding/json"
	"git-codecommit.eu-central-1.amazonaws.com/search-sonata-xq-connector/model"
	fraud "git-codecommit.eu-central-1.amazonaws.com/search-sonata-xq-connector/pStore"
	"log"
)

func HandleMessage(parameterStoreConfig fraud.ParameterStoreConfig, messageBody string) (*TfmResponse, error) {

	// Store timestamp so tracking can be done after message is processed
	//timeIncoming := time.Now()
	var q Query
	err := json.Unmarshal([]byte(messageBody), &q)

	if err != nil {
		log.Println("Couldn't unmarshal messageBody into query:", messageBody, err)
		return nil, err
	}

	//defer awshelper.SendTrackingMessage(awshelper.TrackingMessage{Timestamp: aws.TimeUnixMilli(timeIncoming), ServiceName: awshelper.ServiceName, Action: "Message", SessionId: q.SessionId, RequestId: q.RequestId, ContentType: "JSON", BodyZ: messageBody, Title: "QUEUE_IN", ActionType: "Request"})

	// load config for request
	//_, err = fraud.ServiceConfigForProductChannel(parameterStoreConfig.ConfigProvider, parameterStoreConfig.Servicename, q.Product, q.Channel)

	soapHandler := createSoapHandler()
	tfmResponse, err := soapHandler.Search(q)

	if err != nil {
		log.Println("Something went wrong while searching", err)
		return nil, err
	}

	tfmResponse.Query = handleAntiCorruptionLayer(messageBody)
	/*jsonContent, err := json.Marshal(tfmResponse)
	if err != nil {
		log.Println("Couldn't marshal tfmResponse to json:", tfmResponse, err)
		return "", err
	}*/

	//queueUrl := model.RequestConfig.SqsBaseUrl + model.RequestConfig.QueueNameOut
	// awshelper.SendMessage(queueUrl, string(jsonContent))

	//timeOutgoing := time.Now()
	//defer awshelper.SendTrackingMessage(awshelper.TrackingMessage{Timestamp: aws.TimeUnixMilli(timeOutgoing), ServiceName: awshelper.ServiceName, Action: "Message", SessionId: q.SessionId, RequestId: q.RequestId, ContentType: "JSON", BodyZ: string(jsonContent), Title: "QUEUE: " + queueUrl, ActionType: "Response"})

	return tfmResponse, nil
}

func handleAntiCorruptionLayer(requestJson string) map[string]interface{} {

	log.Println("Putting original query params into tfmResponse..")

	var queryData map[string]interface{}
	err := json.Unmarshal([]byte(requestJson), &queryData)
	if err != nil {
		log.Println("Couldn't unmarshal queryData", err)
	}

	log.Println("Query params successfully put into tfmResponse")

	return queryData
}

func createSoapHandler() SoapHandler {
	return NewSoapHandler(WebserviceConfig{
		WsUrl:            model.RequestConfig.WebserviceUrl,
		WsUser:           model.RequestConfig.WebserviceUser,
		WsPassword:       model.RequestConfig.WebservicePassword,
		Level1:           model.RequestConfig.WebserviceLevel1,
		WsSessionTimeout: model.RequestConfig.WebserviceSessionTimeout,
		ReqModus:         model.RequestConfig.WebserviceRequestType,
		Agent:            model.RequestConfig.WebserviceAgent,
		Instance:         model.RequestConfig.WebserviceInstance,
	})
}
