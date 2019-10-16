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
		Servicename:    "SEARCH-AFKLM-CONNECTOR",
	}

	model.RequestConfig = &model.Config{
		Id:                       "id",
		WebserviceUrl:            "https://ndc-rct.airfranceklm.com/passenger/distribmgmt/001448v01/EXT?",
		WebserviceUser:           "w06014962",
		WebservicePassword:       "TESTPASS",
		WebserviceSessionTimeout: "300",
		WebserviceLevel1:         "NBSO",
		WebserviceAgent:          "1234",
		WebserviceRequestType:    "Live",
		WebserviceInstance:       "NBSO",

		WebserviceParty:                 "MASHERY",
		WebserviceConsumer:              "w06014962",
		WebserviceConsumerLocation:      "External",
		WebserviceConsumerType:          "A",
		WebserviceRelationshipType:      "http://www.af-klm.com/soa/tracking/PrecededBy",
		WebserviceRelatesToValue:        "89562767-4cbkk-4e90-a159-1070b25992fc",
		WebserviceMessageId:             "b762bf9e-2487-42a3-bc88-be998364e51d",
		WebserviceVersion:               "17.1",
		WebserviceCorrelationID:         "3f9ddd-dc6b-41c9-8d4e-8594182ed050",
		WebserviceProperName:            "Test",
		WebservicePseudoCity:            "PARMM211L",
		WebserviceIATA_Nbr:              "12345675",
		WebserviceSequence:              12,
		WebserviceSystemID:              "MAS",
		WebserviceAPIKey:                "qphhwxy2pf5hyeh37gvjdg4c",
		WebserviceParticipantName:       "MASHERY",
		WebserviceCabinCode:             "5",
		WebserviceFarePreferencesCotext: "TO",
		WebserviceFarePreferencesCode:   "758",
		WebserviceFareCode:              "1234",
	}

	tfmResponse, err := search.HandleMessage(config, `{
       "origin" :"AMS",
       "depDate":"2019-09-24",
       "retDate":"2019-09-27",
       "destination":"LHR",
       "adult": 2,
       "child":1,
       "infant":1,
       "journeyType":"ROUNDTRIP",
       "source":"AF"
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
