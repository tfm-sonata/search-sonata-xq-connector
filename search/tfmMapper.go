package search

import (
	"git-codecommit.eu-central-1.amazonaws.com/search-sonata-xq-connector/util"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	routes       map[string]Route
	segments     map[string]Segment
	priceClasses map[string]PriceClass
	passengers   map[string]PassengerDetail
	fareGroupMap map[string]FareGroup
	ancillaries  map[string]Ancillary
	hourRegex    *regexp.Regexp
	minuteRegex  *regexp.Regexp
	dayRegex     *regexp.Regexp
)

func init() {
	hourRegex = regexp.MustCompile(`\d+H`)
	minuteRegex = regexp.MustCompile(`\d+M`)
	dayRegex = regexp.MustCompile(`\d+D`)
}

type TfmMapper interface {
	CreateTFMResponse(rs *AirShoppingRS, conversationToken string, cookie string, responseTime time.Duration, responseId string) (*TfmResponse, error)
	CreateEmptyTfmResponse() *TfmResponse
}

func NewTfmMapper() TfmMapper {
	return &TfmMapperImpl{}
}

type TfmMapperImpl struct {
}

func (*TfmMapperImpl) CreateTFMResponse(rs *AirShoppingRS, conversationToken string, cookie string, responseTime time.Duration, responseId string) (*TfmResponse, error) {

	routes = make(map[string]Route)
	segments = make(map[string]Segment)
	priceClasses = make(map[string]PriceClass)
	passengers = make(map[string]PassengerDetail)
	ancillaries = make(map[string]Ancillary)
	fareGroupMap = make(map[string]FareGroup)

	loadRouteMap(rs, routes)
	loadSegmentMap(rs.DataLists.FlightSegmentList.FlightSegment, segments)
	loadFareGroupMap(rs.DataLists.FareList.FareGroup, fareGroupMap)
	loadPriceClassMap(rs.DataLists.PriceClassList.PriceClass, priceClasses)

	loadPassengerMap(rs.DataLists.PassengerList.Passenger, passengers)
	log.Println("Mapping AirShoppingRS to tfmModel..")

	var combinations = []Combination{}
	createDefaultHandBaggage()

	combinationCounter := 0
	for _, airlineOffer := range rs.OffersGroup.AirlineOffers {

		for _, offer := range airlineOffer.Offer {
			var isValidPriceClassCombination bool = true
			owner := string(*offer.Owner)
			log.Println("Fares for Airline :", owner)
			vcc := "BA"
			//vcc := string(*offer.ValidatingCarrier)

			if combinationCounter < 50 {
				totalFareAmount := offer.TotalPrice.SimpleCurrencyPrice.Value
				var totalTaxAmount float64 = 0
				for _, offerItem := range offer.OfferItem {

					totalTaxAmount = totalTaxAmount + offerItem.FareDetail[0].Price.Taxes.Total.Value
				}
				//create combination
				combination := Combination{}
				additionalParams := make(map[string]string)
				// initialize array
				combination.TotalFareAmount = totalFareAmount
				combination.TotalTaxAmount = totalTaxAmount

				//add route references
				for _, flightservice := range offer.OfferItem[0].Service {
					combination.RouteIDs = append(combination.RouteIDs, string(flightservice.FlightRefs))
				}
				offerItemIdsValues := []string{}
				for _, offerItem := range offer.OfferItem {
					if len(offerItem.FareDetail[0].FareComponent) > 1 {
						isValidPriceClassCombination = validatePriceClassCombination(string(offerItem.FareDetail[0].FareComponent[0].PriceClassRef),
							string(offerItem.FareDetail[0].FareComponent[1].PriceClassRef))
					}
					if isValidPriceClassCombination {
						combination.Fares = createFares(*offerItem, combination.Fares, vcc)
						additionalParamsValue := (string(offerItem.OfferItemID))
						offerItemIdsValues = append(offerItemIdsValues, additionalParamsValue)
						additionalParams[string(offerItem.OfferItemID)] = strings.Trim(string(offerItem.Service[0].PassengerRefs), " ")
					} else {
						log.Println("mismtach in price classes skipping the offer ", string(offer.OfferID))
						break
					}
				}
				if isValidPriceClassCombination {
					offerItemIds := strings.Join(offerItemIdsValues, ",")
					additionalParams["offerId"] = string(offer.OfferID)
					additionalParams["offerItemIds"] = offerItemIds
					//additionalParams["offerValidity"] = offer.OfferExpirationDateTime.Value
					combination.AdditionalParams = additionalParams
					combinationCounter++
					combinations = append(combinations, combination)
				}
			} else {
				break
			}
		}

	}

	additionalParams := make(map[string]string)
	//additionalParams["afklmSearch.conversationToken"] = conversationToken
	//additionalParams["afklmSearch.cookie"] = cookie
	additionalParams["responseId"] = string(*rs.ShoppingResponseID.ResponseID.Value)

	var ancillaryList []Ancillary
	for _, ancillary := range ancillaries {
		ancillaryList = append(ancillaryList, ancillary)
	}

	responseTimes := make(map[string]time.Duration)
	responseTimes["afklmSearch.AirShoppingRQ"] = responseTime

	log.Println("Mapping done")

	return &TfmResponse{
		Combinations:     combinations,
		Routes:           routes,
		Segments:         segments,
		Ancillaries:      ancillaryList,
		AdditionalParams: additionalParams,
		ResponseTimes:    responseTimes,
	}, nil
}
func (*TfmMapperImpl) CreateEmptyTfmResponse() *TfmResponse {
	return &TfmResponse{
		Combinations:     []Combination{},
		Routes:           make(map[string]Route),
		Segments:         make(map[string]Segment),
		Ancillaries:      []Ancillary{},
		AdditionalParams: make(map[string]string),
		ResponseTimes:    make(map[string]time.Duration),
	}
}

func validatePriceClassCombination(outBoundPC string, inBoundPC string) bool {
	outBoundPriceClass := priceClasses[outBoundPC]
	inBoundPriceClass := priceClasses[inBoundPC]
	if outBoundPriceClass.Name == "ECONOMY" || outBoundPriceClass.Name == "Club Europe" ||
		inBoundPriceClass.CabinProduct == "ECONOMY" || inBoundPriceClass.CabinProduct == "Club Europe" {
		return outBoundPriceClass.Name == inBoundPriceClass.Name
	} else {
		return false
	}
}

func loadSegmentMap(segmentList []*ListOfFlightSegmentType, segmentMap map[string]Segment) map[string]Segment {
	for _, segment := range segmentList {
		segment := createSegment(segment)
		segmentMap[segment.Id] = segment
	}

	return segmentMap
}

func loadRouteMap(response *AirShoppingRS, routeMap map[string]Route) map[string]Route {
	for i := 0; i < len(response.DataLists.FlightList.Flight); i++ {

		route := Route{
			Id:                       string(*response.DataLists.FlightList.Flight[i].FlightKey),
			Stops:                    calcuateStops(string(response.DataLists.FlightList.Flight[i].SegmentReferences.Value)),
			SegmentIDs:               strings.Split(string(response.DataLists.FlightList.Flight[i].SegmentReferences.Value), " "),
			ElapsedFlyingTimeMinutes: calculateElapsedFlyingTime(string(response.DataLists.FlightList.Flight[i].Journey.Time)),
		}

		routeMap[route.Id] = route
	}

	return routeMap
}

func loadPriceClassMap(priceClassList []*PriceClassType, priceClassMap map[string]PriceClass) map[string]PriceClass {
	for _, priceClass := range priceClassList {
		priceClass := createPriceClass(priceClass)
		priceClassMap[priceClass.Id] = priceClass
	}

	return priceClassMap
}

func loadPassengerMap(passengerTypeList []*PassengerType, passengerMap map[string]PassengerDetail) map[string]PassengerDetail {
	for _, passengerType := range passengerTypeList {
		passenger := createPassenger(passengerType)
		passengerMap[passenger.Id] = passenger
	}

	return passengerMap
}

func createPassenger(passengerType *PassengerType) PassengerDetail {
	passengerDetail := PassengerDetail{
		Id:   string(passengerType.PassengerID),
		Type: string(*passengerType.PTC), //TODO might have to add a switch case based setting
	}

	return passengerDetail
}

func createPriceClass(priceClassType *PriceClassType) PriceClass {
	var serviceClassList []ServiceClass
	for _, classOfService := range priceClassType.ClassOfService {
		serviceClassList = append(serviceClassList, createServiceClass(*classOfService))
	}

	priceClass := PriceClass{
		Name:             string(*priceClassType.Name),
		Id:               string(priceClassType.PriceClassID),
		ServiceClassList: serviceClassList,
	}

	return priceClass
}

func getCabinProduct(airlineCabinProduct string) string {
	/*var cabinProduct string

	switch airlineCabinProduct {
	case KLM_FLEX:
		cabinProduct = CABIN_ECONOMY
	case AF_ECONOMY_BASIC:
		cabinProduct = CABIN_ECONOMY
	//default:
	//	cabinProduct = CABIN_ECONOMY
	}*/
	return airlineCabinProduct
}
func createServiceClass(classOfService ClassOfService) ServiceClass {

	var fareGroupList []FareGroup
	classOfServiceRefs := string(*classOfService.Refs)
	cosRefItems := strings.Split(classOfServiceRefs, " ")

	for _, cosRefItem := range cosRefItems {
		if strings.Contains(cosRefItem, FARE_BASE_CODE_PREFIX) {
			fareGroupList = append(fareGroupList, fareGroupMap[cosRefItem])
		}
	}

	serviceClass := ServiceClass{
		Id:            classOfServiceRefs,
		Code:          string(*classOfService.Code.Value),
		MarketingName: string(*classOfService.MarketingName.Value),
		FareBaseGroup: fareGroupList,
	}

	return serviceClass
}

func createSegment(segment *ListOfFlightSegmentType) Segment {
	var departureTerminalName = ""
	if segment.Departure.FlightDepartureType != nil && segment.Departure.FlightDepartureType.Terminal != nil {
		departureTerminalName = string(*segment.Departure.FlightDepartureType.Terminal.Name)
	}
	var arrivalTerminalName = ""
	if segment.Arrival != nil && segment.Arrival.Terminal != nil {
		arrivalTerminalName = string(*segment.Arrival.Terminal.Name)
	}
	s := Segment{
		Id:                  string(*segment.SegmentKey),
		Origin:              string(*segment.Departure.AirportCode.Value),
		OriginTerminal:      departureTerminalName,
		Destination:         string(*segment.Arrival.AirportCode.Value),
		DestinationTerminal: arrivalTerminalName,
		FlightNumber:        string(segment.MarketingCarrier.FlightNumber.Value),
		MarketingCarrier:    string(segment.MarketingCarrier.AirlineID.Value),
		OperationCarrier:    string(segment.OperatingCarrier.AirlineID.Value),
	}

	s.DepartureTime = string(segment.Departure.Date) + "T" + string(*segment.Departure.Time)
	s.ArrivalTime = string(segment.Arrival.Date) + "T" + string(*segment.Arrival.Time)

	return s
}

func createDefaultHandBaggage() {
	var ancillary = Ancillary{
		Id:   DEFAULT_HAND_BAGGAGE,
		Type: BAGGAGE_ANCILLARY,
	}
	additionalParams := make(map[string]string)
	additionalParams["peices"] = "1"
	additionalParams["type"] = DEFAULT_HAND_BAGGAGE

	ancillary.AdditionalParams = additionalParams
	ancillaries[DEFAULT_HAND_BAGGAGE] = ancillary
}

func createFares(offerItem OfferItemType, fares []TfmFare, vcc string) []TfmFare {
	passengerRefs := strings.Split(string(offerItem.FareDetail[0].PassengerRefs.Value), " ")
	for _, paxRef := range passengerRefs {
		fares = append(fares, createPassengerFare(offerItem, paxRef, vcc))
	}

	return fares
}

func createPassengerFare(offerItem OfferItemType, paxRef string, vcc string) TfmFare {
	var fareProducts []FareProduct
	tfmFare := TfmFare{
		PaxId:        paxRef,
		PaxType:      passengers[paxRef].Type,
		FareAmount:   offerItem.FareDetail[0].Price.BaseAmount.Value,
		TaxAmount:    util.Round(offerItem.FareDetail[0].Price.Taxes.Total.Value, 100),
		FareProducts: createFareProducts(*offerItem.FareDetail[0], fareProducts, paxRef),
		Vcc:          vcc,
	}

	return tfmFare
}

func createFareProducts(fareDetail FareDetailType, fareProducts []FareProduct, paxRef string) []FareProduct {

	for _, fareComponent := range fareDetail.FareComponent {
		var segmentRefs []string = strings.Split(string(fareComponent.SegmentRefs.Value), " ")

		for _, segmentRef := range segmentRefs {
			fareProducts = append(fareProducts, createFareProduct(*fareComponent, segmentRef, paxRef))
		}
	}

	return fareProducts
}

func createFareProduct(fareComponent FareComponentType, segmentRef string, paxRef string) FareProduct {
	var priceClassRef string = string(fareComponent.PriceClassRef)
	var defaultAncillaries []string
	var fareProductFinal FareProduct
	defaultAncillaries = append(defaultAncillaries, ancillaries[DEFAULT_HAND_BAGGAGE].Id)
	var priceClass = priceClasses[priceClassRef] //PCR_1

	for _, classOfServiceEntry := range priceClass.ServiceClassList {
		if strings.Contains(classOfServiceEntry.Id, segmentRef) {
			fareList := strings.Split(classOfServiceEntry.Id, " ") ///////////////////rename////////////////////////////////////////////////////////rename

			for _, fareitem := range fareList {

				if strings.Contains(fareitem, FARE_BASE_CODE_PREFIX) {
					var fareGroup = fareGroupMap[fareitem]
					fareProductFinal = FareProduct{
						SegmentID:    segmentRef,
						CabinProduct: classOfServiceEntry.MarketingName,
						FareBase:     fareGroup.FareBasisCode,
						AncillaryIDs: defaultAncillaries,
					}
				}

			}
		}
	}

	return fareProductFinal
}

func loadFareGroupMap(fareGroupList []*FareGroup__1, fareGroupMap map[string]FareGroup) map[string]FareGroup {
	for _, fareGroupListItem := range fareGroupList {
		fareGroup := createFareGroup(fareGroupListItem)
		fareGroupMap[fareGroup.Id] = fareGroup
	}

	return fareGroupMap
}
func createFareGroup(fareGroupListItem *FareGroup__1) FareGroup {
	fareGroup := FareGroup{
		Id:            fareGroupListItem.ListKey,
		FareBasisCode: string(*fareGroupListItem.FareBasisCode.Code),
		FareCode:      string(*fareGroupListItem.Fare.FareCode),
	}

	return fareGroup
}

func calcuateStops(segmentReferences string) int8 {
	return int8(len(strings.Split(segmentReferences, " ")) - 1)
}

func getTime(flightDuration string) []string {
	return strings.Split(flightDuration, "T")
}

func calculateElapsedFlyingTime(flightDuration string) int {
	var elapsedFlyingTimeMinutes int

	flightDurationDayTime := getTime(flightDuration)
	log.Println(flightDurationDayTime)

	dayDuration := flightDurationDayTime[0]
	timeDuration := flightDurationDayTime[1]

	log.Println(dayDuration + " " + timeDuration)

	dayString := dayRegex.FindString(dayDuration)

	if dayString != "" {
		dayString := strings.TrimSuffix(dayString, "D")
		days, err := strconv.Atoi(dayString)
		if err != nil {
			log.Println("Couldn't elapsedFlyingTimeMinutes", flightDuration, err)
		}
		elapsedFlyingTimeMinutes = elapsedFlyingTimeMinutes + days*24*60
	}

	hourString := hourRegex.FindString(timeDuration)

	if hourString != "" {
		hourString := strings.TrimSuffix(hourString, "H")
		hours, err := strconv.Atoi(hourString)
		if err != nil {
			log.Println("Couldn't elapsedFlyingTimeMinutes", flightDuration, err)
		}
		elapsedFlyingTimeMinutes = elapsedFlyingTimeMinutes + hours*60
	}

	minuteString := minuteRegex.FindString(timeDuration)
	if minuteString != "" {
		minuteString := strings.TrimSuffix(minuteString, "M")
		minutes, err := strconv.Atoi(minuteString)
		if err != nil {
			log.Println("Couldn't elapsedFlyingTimeMinutes", flightDuration, err)
		}
		elapsedFlyingTimeMinutes = elapsedFlyingTimeMinutes + minutes
	}

	return elapsedFlyingTimeMinutes
}
