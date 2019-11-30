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
	//var cabinCode CodesetValueSimpleType = CodesetValueSimpleType(model.RequestConfig.WebserviceCabinCode) //5
	//var farePreferenceceType := model.RequestConfig.WebserviceFareCode
	var passengers []*PassengerType = createPaxTypeList(query)
	var originDestinations *AirShopReqAttributeQueryType = craeteOriginDestinations(query)
	//var farePreferences *FarePreferences = createFarePreference(query) //nil

	pointOfSale, _ := createPointOfSale(query)
	document, _ := createDocument(query)
	party, _ := createParty(query)

	airShoppingRQTemp := &AirShoppingRQ{
		Version:     airShoppingRQVersion,
		PointOfSale: pointOfSale,
		Document:    document,
		Party:       party,
		CoreQuery: &CoreQuery__1{
			OriginDestinations: originDestinations,
		},
		/*Preference: &Preference__1{
			CabinPreferences: &CabinPreferences{
				CabinType: []*CabinType{{
					CodesetType: &CodesetType{
						Code: &cabinCode,
					},
				}},
			},
			FarePreferences: farePreferences,
		},*/
		DataLists: &DataLists__1{
			PassengerList: &PassengerList__1{Passenger: passengers},
		},
	}
	return airShoppingRQTemp, nil
}

func createPointOfSale(query Query) (*PointOfSale, error) {
	var requestTime = time.Now().Format(DATE_FORMAT)
	var countryCode CountrySimpleType = CountrySimpleType(model.RequestConfig.WebServiceRequestCountryCode)
	var cityCode AirportCitySimpleType = AirportCitySimpleType(model.RequestConfig.WebServiceRequestCityCode)

	location := &Location__1{
		CountryCode: &CountryCode{CountryCodeType: &CountryCodeType{Value: countryCode}},
		CityCode:    &CityCode{CityCodeType: &CityCodeType{Value: cityCode}},
	}
	pointOfSale := &PointOfSale{RequestTime: &RequestTime__1{Value: requestTime},
		Location: location}
	return pointOfSale, nil
}

func createDocument(query Query) (*Document, error) {
	var name ProperNameSimpleType = ProperNameSimpleType(model.RequestConfig.WebServiceDocumentName)
	var refVersion ContextSimpleType = ContextSimpleType(model.RequestConfig.WebserviceVersion)
	return &Document{
		Name:             &name,
		ReferenceVersion: &refVersion,
	}, nil
}

func createAgentUserSender(userId string) *AgentUserSender {
	agentUserName := ProperNameSimpleType(model.RequestConfig.WebServiceSenderName)
	agencyCoreRepType := &AgencyCoreRepType{SellerCoreRepType: &SellerCoreRepType{
		Name: &agentUserName,
	},
	}

	agentUserId := AgentUserID__1{
		UniqueIDContextType: &UniqueIDContextType{Value: UniqueStringID_SimpleType(userId)},
	}
	agentUserType := AgentUserType{
		AgencyCoreRepType: agencyCoreRepType,
		AgentUserID:       &agentUserId,
		UserRole:          nil,
	}

	agentUserSender := AgentUserSender{
		AgentUserMsgPartyCoreType: &AgentUserMsgPartyCoreType{AgentUserType: &agentUserType},
	}

	return &agentUserSender
}

func createParty(query Query) (*Party, error) {
	//var properNameSimpleType ProperNameSimpleType = ProperNameSimpleType(model.RequestConfig.WebserviceProperName)
	//var pseudoCitySimpleType PseudoCitySimpleType = PseudoCitySimpleType(model.RequestConfig.WebservicePseudoCity)
	//var iATA_NbrSimpleType IATA_NbrSimpleType = IATA_NbrSimpleType(model.RequestConfig.WebserviceIATA_Nbr)
	//var id UniqueStringID_SimpleType = UniqueStringID_SimpleType(model.RequestConfig.Id)
	//var agentUserID UniqueStringID_SimpleType = UniqueStringID_SimpleType(model.RequestConfig.WebserviceAgent)
	/*var sequenceNumber int32 = int32(model.RequestConfig.WebserviceSequence)
	var participantName ProperNameSimpleType = ProperNameSimpleType(model.RequestConfig.WebserviceParticipantName)
	var systemID UniqueStringID_SimpleType = UniqueStringID_SimpleType(model.RequestConfig.WebserviceSystemID)*/
	//var name ProperNameSimpleType = ProperNameSimpleType("TUI")
	//var aggrId UniqueStringID_SimpleType = UniqueStringID_SimpleType("00038566")

	var airLineID AirlineDesigSimpleType = AirlineDesigSimpleType(query.Source)
	//var airLineName CarrierNameType = CarrierNameType(SourceMap[query.Source])
	var recipientName = CarrierNameType(model.RequestConfig.WebserviceAirLineName)
	userName := model.RequestConfig.WebserviceParticipantName
	airlineCoreType := &AirlineCoreRepType{
		AirlineID: &AirlineID{
			AirlineID_Type: &AirlineID_Type{Value: airLineID},
		},
		Name: &recipientName,
	}
	party := &Party{
		Sender: &Sender__1{
			//AgentUserSender: createAgentUserSender(userName),
			AgentUserSender: createAgentUserSender(userName),
		},
		Recipient: &Recipient__1{
			ORA_Recipient: &ORA_Recipient{
				AirlineMsgPartyCoreType: &AirlineMsgPartyCoreType{
					AirlineCoreRepType: airlineCoreType,
				},
			},
		},
		/*TravelAgencySender: &TravelAgencySender{
		TrvlAgencyMsgPartyCoreType: &TrvlAgencyMsgPartyCoreType{
			TravelAgencyType: &TravelAgencyType{
				AgencyCoreRepType: &AgencyCoreRepType{
		*/
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
		/*Participants: &Participants__1{Participant: []*Participant__1{
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
		}},*/
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

func createPassengerType(paxNo int8, paxType PassengerTypeCodeType, query Query) *PassengerType {
	var paxId = string(query.Source + "_" + PAX + "_" + strconv.Itoa(int(paxNo+1)))
	return &PassengerType{
		PassengerID: paxId,
		PTC:         &paxType,
	}

}
func createPaxTypeList(query Query) []*PassengerType {
	var paxTypeList []*PassengerType
	var adt PassengerTypeCodeType = PassengerTypeCodeType(ADT)
	var chd PassengerTypeCodeType = PassengerTypeCodeType(CHD)
	var inf PassengerTypeCodeType = PassengerTypeCodeType(INF)

	totalPaxCount := query.Adult + query.Child + query.Infant

	var i, adtCount, chdCount, infCount int8

	for i < totalPaxCount {
		for ; adtCount < query.Adult; adtCount++ {
			paxTypeList = append(paxTypeList, createPassengerType(i, adt, query))
			i++
		}
		for ; chdCount < query.Child; chdCount++ {
			paxTypeList = append(paxTypeList, createPassengerType(i, chd, query))
			i++
		}
		for ; infCount < query.Infant; infCount++ {
			paxTypeList = append(paxTypeList, createPassengerType(i, inf, query))
			i++
		}

	}

	return paxTypeList
}

func craeteOriginDestinations(query Query) *AirShopReqAttributeQueryType {

	var originDestinations *AirShopReqAttributeQueryType
	var departureAirportCode AirportCitySimpleType = AirportCitySimpleType(query.Origin) //"AMS"
	var departureDate string = string(query.DepDate)
	var returnDate string = string(query.RetDate)                                           //"2019-12-15Z"
	var arrivalAirportCode AirportCitySimpleType = AirportCitySimpleType(query.Destination) //"LHR"
	/*var calendarDate CalendarDates__1 = CalendarDates__1{
		DaysBefore: 0,
		DaysAfter:  0,
	}*/

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
					//CalendarDates: &calendarDate,

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
					//CalendarDates: &calendarDate,
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
