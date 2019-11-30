package search

import "time"

type Query struct {
	Origin          string `json:"origin"`
	Destination     string `json:"destination"`
	DepDate         string `json:"depDate"`
	RetDate         string `json:"retDate"`
	Currency        string `json:"currency"`
	SessionId       string `json:"sessionId"`
	RequestId       string `json:"requestId"`
	Product         string `json:"product"`
	Channel         string `json:"channel"`
	Adult           int8   `json:"adult"`
	Child           int8   `json:"child"`
	Infant          int8   `json:"infant"`
	JourneyType     string `json:"journeyType"`
	TrackingEnabled bool   `json:"trackingEnabled"`
	Source          string `json:"source"`
}

type Route struct {
	Id                       string            `json:"id"`
	Stops                    int8              `json:"stops"`
	ElapsedFlyingTimeMinutes int               `json:"elapsedFlyingTimeMinutes"`
	SegmentIDs               []string          `json:"segmentIDs"`
	AdditionalParams         map[string]string `json:"additionalParams,omitempty"`
}

type Segment struct {
	Id                  string            `json:"id"`
	Origin              string            `json:"origin"`
	OriginTerminal      string            `json:"originTerminal"`
	Destination         string            `json:"destination"`
	DestinationTerminal string            `json:"destinationTerminal"`
	DepartureTime       string            `json:"departureTime"`
	ArrivalTime         string            `json:"arrivalTime"`
	FlightNumber        string            `json:"flightNumber"`
	AirplaneType        string            `json:"airplaneType"`
	MarketingCarrier    string            `json:"marketingCarrier"`
	OperationCarrier    string            `json:"operatingCarrier"`
	AdditionalParams    map[string]string `json:"additionalParams,omitempty"`
}

type Combination struct {
	TotalFareAmount  float64           `json:"totalFareAmount"`
	TotalTaxAmount   float64           `json:"totalTaxAmount"`
	Fares            []TfmFare         `json:"fares"`
	RouteIDs         []string          `json:"routeIDs"`
	TariffType       string            `json:"tariffType"`
	AdditionalParams map[string]string `json:"additionalParams,omitempty"`
}

type TfmFare struct {
	PaxId        string        `json:"paxId"`
	PaxType      string        `json:"paxType"`
	FareAmount   float64       `json:"fareAmount"`
	TaxAmount    float64       `json:"taxAmount"`
	FareProducts []FareProduct `json:"fareProducts"`
	Vcc          string        `json:"vcc"`
}

type FareProduct struct {
	SegmentID        string            `json:"segmentID"`
	CabinProduct     string            `json:"cabinProduct"`
	FareBase         string            `json:"fareBase"`
	AncillaryIDs     []string          `json:"ancillaryIDs"`
	AdditionalParams map[string]string `json:"additionalParams,omitempty"`
}

type Ancillary struct {
	Id               string            `json:"id"`
	Type             string            `json:"type"`
	AdditionalParams map[string]string `json:"additionalParams,omitempty"`
}

type TfmResponse struct {
	Query            map[string]interface{}   `json:"query"`
	Routes           map[string]Route         `json:"routes"`
	Segments         map[string]Segment       `json:"segments"`
	Combinations     []Combination            `json:"combinations"`
	Ancillaries      []Ancillary              `json:"ancillaries"`
	AdditionalParams map[string]string        `json:"additionalParams,omitempty"`
	ResponseTimes    map[string]time.Duration `json:"responseTimes,omitempty"`
}

type ReservationClassHolder struct {
	BookingClass string
	FareBase     string
	CabinProduct string
}

type PriceClass struct {
	Name             string
	Id               string //PriceClassID=PC1
	FareBase         string
	CabinProduct     string
	ServiceClassList []ServiceClass
}

type ServiceClass struct {
	Id            string
	Code          string
	MarketingName string
	FareBaseGroup []FareGroup
}

type FareGroup struct {
	Id            string
	FareBasisCode string
	FareCode      string
}

type FareComponentAugPointMetadata struct {
	BaseFare        string
	BaseFareCode    string
	DisplayFare     string
	DisplayFareCode string
	Discount        string
	DiscountCode    string
	MetaDataKeyRef  string
}
type FareDetailAugPointMetadata struct {
	FareType       string
	FareLevel      string
	FareId         string
	MetaDataKeyRef string
}

type PassengerDetail struct {
	Id   string //PAX1
	Type string
}

const (
	NBSO_TO   = "NBSO_TO"
	NBSO_PUB  = "NBSO_PUB"
	ADT       = "ADT"
	CHD       = "CHD"
	INF       = "INF"
	PAX       = "PAX_"
	TARIFF_TO = "TO"
	TARIFF_PU = "PU"

	CABIN_ECONOMY           = "ECONOMY"
	KLM_FLEX                = "Flex"
	KLM_STANDARD            = "Standard"
	KLM_LIGHT               = "Light"
	AF_ECONOMY_SMART        = "ECONOMY SMART"
	AF_ECONOMY_BASIC        = "ECONOMY BASIC"
	AF_ECONOMY_BASICPLUS    = "ECONOMY BASICPLUS"
	ROUNDTRIP               = "ROUNDTRIP"
	ONEWAY                  = "ONEWAY"
	DEFAULT_HAND_BAGGAGE    = "HAND"
	DEFAULT_CHECKED_BAGGAGE = "BAGGAGE"
	BAGGAGE_ANCILLARY       = "BAGGAGE"
	BAGGAGE_NORMAL          = "NORMAL"
	DATE_FORMAT             = "2006-01-02T15:04:05"
	FARE_BASE_CODE_PREFIX   = "FBCODE"

	FARE_TYPE_SUFFIX  = "_FareType"
	FARE_LEVEL_SUFFIX = "_FareLevel"
	FARE_ID_SUFFIX    = "_FareId"

	BASE_FARE_SUFFIX         = "_BaseFare"
	BASE_FARE_CODE_SUFFIX    = "_BaseFareCode"
	DISPLAY_FARE_SUFFIX      = "_DisplayFare"
	DISPLAY_FARE_CODE_SUFFIX = "_DisplayFareCode"
	DISC_FARE_SUFFIX         = "_Discount"
	DISC_FARE_CODE_SUFFIX    = "_DiscountCode"

	METADATA_KEY_SUFFIX = "_MetaDataKey"

	SPACE_DELIMITER = " "
	COMMA_DELIMITER = ","

	OD_REF = "OdRef"
)

var SourceMap = map[string]string{
	"AF": "AIRFRANCE",
	"KL": "KLM",
}
