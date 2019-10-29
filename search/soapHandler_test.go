package search_test

import (
	"encoding/json"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/search-sonata-xq-connector/model"
	"git-codecommit.eu-central-1.amazonaws.com/search-sonata-xq-connector/search"
	"log"
	"testing"
)

var (
	q      search.Query
	cookie string
)

func TestSoapHandler(t *testing.T) {
	//t.Skip()

	model.RequestConfig = &model.Config{
		Id:                       "id",
		WebserviceUrl:            "https://iflyrestest.ibsgen.com:6013/iRes_NdcRes_WS/services/NdcResService172SOAPPort",
		WebserviceUser:           "jetradar",
		WebservicePassword:       "oPBCDECGZCRA",
		WebserviceSessionTimeout: "300",
		WebserviceLevel1:         "NBSO",
		WebserviceAgent:          "1234",
		WebserviceRequestType:    "Live",
		WebserviceInstance:       "NBSO",

		WebserviceParty:                 "jetradar",
		WebserviceConsumer:              "jetradar",
		WebserviceConsumerLocation:      "External",
		WebserviceConsumerType:          "A",
		WebserviceRelationshipType:      "",
		WebserviceRelatesToValue:        "89562767-4cbkk-4e90-a159-1070b25992fc",
		WebserviceMessageId:             "b762bf9e-2487-42a3-bc88-be998364e51d",
		WebserviceVersion:               "17.2",
		WebserviceCorrelationID:         "3f9ddd-dc6b-41c9-8d4e-8594182ed050",
		WebserviceProperName:            "Test",
		WebservicePseudoCity:            "PARMM211L",
		WebserviceIATA_Nbr:              "12345675",
		WebserviceSequence:              12,
		WebserviceSystemID:              "MAS",
		WebserviceAPIKey:                "qphhwxy2pf5hyeh37gvjdg4c",
		WebserviceParticipantName:       "JETRADAR",
		WebserviceCabinCode:             "5",
		WebserviceFarePreferencesCotext: "TO",
		WebserviceFarePreferencesCode:   "758",
		WebserviceFareCode:              "1234",
		WebServiceDocumentName:          "NDC",
		WebServiceRequestCountryCode:    "TR",
		WebServiceRequestCityCode:       "City",
		WebserviceAirLineID:             "XQ",
		WebserviceAirLineName:           "Sun Express",
		WebServiceContentType:           "application/xml",
		WebServiceSenderName:            "Guest Website",
	}

	soapHandler := search.NewSoapHandler(search.WebserviceConfig{
		WsUrl:            model.RequestConfig.WebserviceUrl,
		WsUser:           model.RequestConfig.WebserviceUser,
		WsPassword:       model.RequestConfig.WebservicePassword,
		Level1:           model.RequestConfig.WebserviceLevel1,
		WsSessionTimeout: model.RequestConfig.WebserviceSessionTimeout,
		ReqModus:         model.RequestConfig.WebserviceRequestType,
		Agent:            model.RequestConfig.WebserviceAgent,
		Instance:         model.RequestConfig.WebserviceInstance,
	}).(search.SoapHandler)

	var messageBody = `{
       "origin" :"AYT",
       "depDate":"2019-11-23",
       "retDate":"2019-11-28",
       "destination":"FRA",
       "adult": 1,
       "child":0,
       "infant":0,
       "journeyType":"ROUNDTRIP",
       "source":"XQ"
        }`

	var q search.Query
	err := json.Unmarshal([]byte(messageBody), &q)

	if err != nil {
		log.Println("Couldn't unmarshal messageBody into query:", messageBody, err)
	}

	tfmResponse, err := soapHandler.Search(q)

	if err != nil {
		log.Println("Something went wrong while searching", err)
	}

	searchResponse, error := json.Marshal(tfmResponse)
	if error != nil {
		log.Println("unable to marsha Tfm response", error)
	}

	fmt.Println(string(searchResponse))

}
