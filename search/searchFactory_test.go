package search_test

import (
	"encoding/json"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/search-sonata-xq-connector/model"
	"git-codecommit.eu-central-1.amazonaws.com/search-sonata-xq-connector/search"
	"log"
	"testing"
)

func TestCreateAirShoppingRQ(t *testing.T) {
	//t.Skip()
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

		WebserviceParty:            "MASHERY",
		WebserviceConsumer:         "w06014962",
		WebserviceConsumerLocation: "External",
		WebserviceConsumerType:     "A",
		WebserviceRelationshipType: "http://www.af-klm.com/soa/tracking/PrecededBy",
		WebserviceRelatesToValue:   "89562767-4cbkk-4e90-a159-1070b25992fc",
		WebserviceMessageId:        "b762bf9e-2487-42a3-bc88-be998364e51d",
		WebserviceVersion:          "17.1",
		WebserviceCorrelationID:    "3f9ddd-dc6b-41c9-8d4e-8594182ed050",
		WebserviceProperName:       "Test",
		WebservicePseudoCity:       "PARMM211L",
		WebserviceIATA_Nbr:         "12345675",
		WebserviceSequence:         12,
		WebserviceSystemID:         "MAS",
		WebserviceAirLineID:        "AF",
		WebserviceAirLineName:      "AIRFRANCE",
		WebserviceAPIKey:           "zzk78msv4j9r57r6kkn6hv4g",
		WebserviceParticipantName:  "MASHERY",
	}
	var factories search.SearchFactory = &search.SearchFactoryImpl{}
	query := search.Query{
		Origin:          "AMS",
		Destination:     "LHR",
		DepDate:         "2019-09-13",
		RetDate:         "2019-09-14",
		Currency:        "",
		SessionId:       "47d0950e-b9bb-11e9-9234-0123456789ab",
		RequestId:       "403115e7-344d-4cac-99e3-de0c1141ee93",
		Product:         "TFM",
		Channel:         "WEB",
		Adult:           0,
		Child:           0,
		Infant:          0,
		JourneyType:     "ROUNDTRIP",
		TrackingEnabled: false,
	}

	airShoppingRQ, err := factories.CreateAirShoppingRQ(query)
	if err != nil {
		log.Fatal(err)
	}

	jsonContent, err := json.Marshal(airShoppingRQ)
	fmt.Println("Length of jsonContent: ", len(string(jsonContent)))
	fmt.Println(string(jsonContent))

}
