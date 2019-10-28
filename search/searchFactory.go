package search

import (
	"git-codecommit.eu-central-1.amazonaws.com/search-sonata-xq-connector/model"
	"strconv"
	"time"
)

type SearchFactory interface {
	CreateAirShoppingRQ(query Query) (*AirShoppingRQ, error)
	createSOAPHeader(config *model.Config) (interface{}, error)
}
type SearchFactoryImpl struct {
}

func (*SearchFactoryImpl) createSOAPHeader(config *model.Config) (interface{}, error) {
	var soapHeader = struct {
		AirShopingHeader `xml:"tns:AirShopingHeader"`
	}{
		AirShopingHeader: AirShopingHeader{},
	}
	return soapHeader, nil
}
func (*SearchFactoryImpl) CreateAirShoppingRQ(query Query) (*AirShoppingRQ, error) {

	var airShoppingRQVersion = (model.RequestConfig.WebserviceVersion)
	//var primanryLangID ="EN"
	var cabinCode CodesetValueSimpleType = CodesetValueSimpleType(model.RequestConfig.WebserviceCabinCode) //5
	//var farePreferenceceType := model.RequestConfig.WebserviceFareCode
	var passengers []*PassengerType = createPassangers(query)
	var originDestinations *AirShopReqAttributeQueryType = craeteOriginDestinations(query)
	var farePreferences *FarePreferences = createFarePreference(query) //nil

	//pointOfSale, _ := createPointOfSale(query)
	document, _ := createDocument(query)
	party, _ := createParty(query)

	airShoppingRQTemp := &AirShoppingRQ{
		Version:  airShoppingRQVersion,
		Document: document,
		Party:    party,
		CoreQuery: &CoreQuery__1{
			OriginDestinations: originDestinations,
		},
		Preference: &Preference__1{
			CabinPreferences: &CabinPreferences{
				CabinType: []*CabinType{{
					CodesetType: &CodesetType{
						Code: &cabinCode,
					},
				}},
			},
			FarePreferences: farePreferences,
		},
		DataLists: &DataLists__1{
			PassengerList: &PassengerList__1{Passenger: passengers},
		},
	}
	return airShoppingRQTemp, nil
}

func createPointOfSale(query Query) (*PointOfSale, error) {
	var requestTime = time.Now().Format(DATE_FORMAT)
	pointOfSale := &PointOfSale{RequestTime: &RequestTime__1{Value: requestTime}}
	return pointOfSale, nil
}

func createDocument(query Query) (*Document, error) {
	var name ProperNameSimpleType = ProperNameSimpleType(query.Source)
	return &Document{
		Name: &name,
	}, nil
}

func createParty(query Query) (*Party, error) {
	//var properNameSimpleType ProperNameSimpleType = ProperNameSimpleType(model.RequestConfig.WebserviceProperName)
	//var pseudoCitySimpleType PseudoCitySimpleType = PseudoCitySimpleType(model.RequestConfig.WebservicePseudoCity)
	var iATA_NbrSimpleType IATA_NbrSimpleType = IATA_NbrSimpleType(model.RequestConfig.WebserviceIATA_Nbr)
	var id UniqueStringID_SimpleType = UniqueStringID_SimpleType(model.RequestConfig.Id)
	//var agentUserID UniqueStringID_SimpleType = UniqueStringID_SimpleType(model.RequestConfig.WebserviceAgent)
	/*var sequenceNumber int32 = int32(model.RequestConfig.WebserviceSequence)
	var participantName ProperNameSimpleType = ProperNameSimpleType(model.RequestConfig.WebserviceParticipantName)
	var systemID UniqueStringID_SimpleType = UniqueStringID_SimpleType(model.RequestConfig.WebserviceSystemID)*/
	var name ProperNameSimpleType = ProperNameSimpleType("TUI")
	var aggrId UniqueStringID_SimpleType = UniqueStringID_SimpleType("00038566")

	//var airLineID AirlineDesigSimpleType = AirlineDesigSimpleType(query.Source)
	//var airLineName CarrierNameType = CarrierNameType(SourceMap[query.Source])

	party := &Party{
		Sender: &Sender__1{
			TravelAgencySender: &TravelAgencySender{
				TrvlAgencyMsgPartyCoreType: &TrvlAgencyMsgPartyCoreType{
					TravelAgencyType: &TravelAgencyType{
						AgencyCoreRepType: &AgencyCoreRepType{
							//SellerCoreRepType: &SellerCoreRepType{
							//	Name: &properNameSimpleType,
							//},
							//PseudoCity: &PseudoCity__1{
							//	Value: &pseudoCitySimpleType,
							//},
							IATA_Number: &iATA_NbrSimpleType,
						},
						AgencyID: &AgencyID{
							UniqueIDContextType: &UniqueIDContextType{
								Value: id,
							},
						},
					},
					//AgentUser: &AgentUser{
					//	AgentUserID: &AgentUserID__1{
					//		UniqueIDContextType: &UniqueIDContextType{
					//			Value: agentUserID,
					//		}},
					//},
				},
			},
		},
		/*Participants: &Participants__1{Participant: []*Participant__1{{
			EnabledSystemParticipant: &EnabledSystemParticipant{
				SequenceNumber: sequenceNumber,
				EnabledSysMsgPartyCoreType: &EnabledSysMsgPartyCoreType{
					&EnabledSystemType{
						IntermediaryCoreRepType: &IntermediaryCoreRepType{
							Name: &participantName,
						},
						SystemID: &SystemID{
							UniqueIDContextType: &UniqueIDContextType{
								Value: systemID,
							},
						},
					},
				},
			},
		}}},*/
		Participants: &Participants__1{Participant: []*Participant__1{
			{
				AggregatorParticipant: &AggregatorParticipant{
					AggregatorMsgPartyCoreType: &AggregatorMsgPartyCoreType{
						AggregatorType: &AggregatorType{
							IntermediaryCoreRepType: &IntermediaryCoreRepType{
								Name: &name,
							},
							AggregatorID: &AggregatorID{
								UniqueIDContextType: &UniqueIDContextType{
									Value: aggrId,
								},
							},
						},
					},
					SequenceNumber: 1,
				},
			},
		}},
		//Participants: &Participants__1{Participant: []*Participant__1{{
		//	AggregatorParticipant: &AggregatorParticipant{
		//		SequenceNumber: sequenceNumber,
		//		AggregatorMsgPartyCoreType: &AggregatorMsgPartyCoreType{
		//			&AggregatorType{
		//				IntermediaryCoreRepType: &IntermediaryCoreRepType{
		//					Name: &participantName,
		//				},
		//				AggregatorID: &AggregatorID_Type{
		//					UniqueIDContextType: &UniqueIDContextType{
		//						Value: systemID,
		//					},
		//				},
		//			},
		//		},
		//	},
		//}}}
		//Recipient: &Recipient__1{
		//	ORA_Recipient: &ORA_Recipient{
		//		AirlineMsgPartyCoreType: &AirlineMsgPartyCoreType{
		//			AirlineCoreRepType: &AirlineCoreRepType{
		//				AirlineID: &AirlineID{
		//					AirlineID_Type: &AirlineID_Type{
		//						Value: airLineID,
		//					},
		//				},
		//				Name: &airLineName,
		//			},
		//		},
		//	},
		//},
	}
	return party, nil
}

func createParameters(airShoppingRQ *AirShoppingRQ, query Query) (*AirShoppingRQ, error) {
	return nil, nil
}
func createPassangers(query Query) []*PassengerType {
	passangers := []*PassengerType{}
	var adt PassengerTypeCodeType = PassengerTypeCodeType(ADT)
	var chd PassengerTypeCodeType = PassengerTypeCodeType(CHD)
	var ift PassengerTypeCodeType = PassengerTypeCodeType(INF)

	var i int8
	for i = 0; i < query.Adult; i++ {
		var passengerId = string(PAX + strconv.Itoa(int(i+1))) // (b + strconv.Itoa(a))
		passangers = append(passangers, &PassengerType{
			PassengerID: passengerId,
			PTC:         &adt,
		})
	}
	for i = 0; i < query.Child; i++ {
		var passengerId = string(PAX + strconv.Itoa(int(i+1)))
		passangers = append(passangers, &PassengerType{
			PassengerID: passengerId,
			PTC:         &chd,
		})
	}
	for i = 0; i < query.Infant; i++ {
		var passengerId = string(PAX + strconv.Itoa(int(i+1)))
		passangers = append(passangers, &PassengerType{
			PassengerID: passengerId,
			PTC:         &ift,
		})
	}
	return passangers
}

func craeteOriginDestinations(query Query) *AirShopReqAttributeQueryType {

	var originDestinations *AirShopReqAttributeQueryType
	var departureAirportCode AirportCitySimpleType = AirportCitySimpleType(query.Origin) //"AMS"
	var departureDate string = string(query.DepDate)
	var returnDate string = string(query.RetDate)                                           //"2019-12-15Z"
	var arrivalAirportCode AirportCitySimpleType = AirportCitySimpleType(query.Destination) //"LHR"

	if query.JourneyType == ROUNDTRIP {
		originDestinations = &AirShopReqAttributeQueryType{
			OriginDestination: []*OriginDestination__1{
				{
					Departure: &Departure{
						FlightDepartureType: &FlightDepartureType{
							AirportCode: &AirportCode__1{
								Value: &departureAirportCode,
							},
							Date:        departureDate,
							Time:        nil,
							AirportName: nil,
						},
					},
					Arrival: &Arrival{
						AirportCode: &AirportCode__2{
							Value: &arrivalAirportCode,
						},
					},
				},
				{
					Departure: &Departure{
						FlightDepartureType: &FlightDepartureType{
							AirportCode: &AirportCode__1{
								Value: &arrivalAirportCode,
							},
							Date: returnDate,
						},
					},
					Arrival: &Arrival{
						AirportCode: &AirportCode__2{
							Value: &departureAirportCode,
						},
					},
				},
			},
		}
	} else {
		originDestinations = &AirShopReqAttributeQueryType{
			OriginDestination: []*OriginDestination__1{
				{
					Departure: &Departure{
						FlightDepartureType: &FlightDepartureType{
							AirportCode: &AirportCode__1{
								Value: &departureAirportCode,
							},
							Date:        departureDate,
							Time:        nil,
							AirportName: nil,
						},
					},
					Arrival: &Arrival{
						AirportCode: &AirportCode__2{
							Value: &arrivalAirportCode,
						},
					},
				},
			},
		}
	}

	return originDestinations
}

func createFarePreference(query Query) *FarePreferences {
	//var fareCode CodesetValueSimpleType = CodesetValueSimpleType(model.RequestConfig.WebserviceFareCode)
	var farePreferenceContext ContextSimpleType = ContextSimpleType(model.RequestConfig.WebserviceFarePreferencesCotext)
	var farePreferenceCode IATA_CodeType = IATA_CodeType(model.RequestConfig.WebserviceFarePreferencesCode)
	farePreference := &FarePreferences{
		Types: &Types__1{Type: []*Type__4{
			{Value: &farePreferenceCode, PreferencesContext: &farePreferenceContext},
		}},
		FareCodes: &FareCodes__1{
			/*Code: []*Code__6{{
				FareBasisCodeType: &FareBasisCodeType{Code: &fareCode}},
			},*/
		},
	}

	return farePreference
}
