package search_test

import (
	"encoding/json"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/search-sonata-xq-connector/model"
	"git-codecommit.eu-central-1.amazonaws.com/search-sonata-xq-connector/pStore"
	"git-codecommit.eu-central-1.amazonaws.com/search-sonata-xq-connector/search"
	"log"
	"testing"
)

func init() {
}

func TestHandleMessage(t *testing.T) {
	//t.Skip()
	config := fraud.ParameterStoreConfig{
		ConfigProvider: "common-config-service",
		Servicename:    "SEARCH-XQ-CONNECTOR",
	}

	model.RequestConfig = &model.Config{
		Id:                       "id",
		WebserviceUrl:            "https://iflyrestest.ibsgen.com:6013/iRes_NdcRes_WS/services/NdcResService172SOAPPort?",
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
		WebserviceFarePreferencesCotext: "PU",
		WebserviceFarePreferencesCode:   "758",
		WebserviceFareCode:              "1234",
		WebServiceDocumentName:          "NDC",
		WebServiceRequestCountryCode:    "DE",
		WebServiceRequestCityCode:       "City",
		WebserviceAirLineID:             "XQ",
		WebserviceAirLineName:           "Sun Express",
		WebServiceContentType:           "application/xml",
		WebServiceSenderName:            "Guest Website",
	}

	tfmResponse, err := search.HandleMessage(config, `{
       "origin" :"AYT",
       "depDate":"2019-12-29",
       "retDate":"2019-12-30",
       "destination":"FRA",
       "adult": 2,
       "child":1,
       "infant":1,
       "journeyType":"ROUNDTRIP",
       "source":"XQ"
        }`)
	if err != nil {
		t.Errorf("createSession() failed.")
	}
	//fmt.Println(tfmResponse)
	searchResponse, error := json.Marshal(tfmResponse)
	if error != nil {
		log.Println("unable to marsha Tfm response", error)
	}

	fmt.Println(string(searchResponse))

}
