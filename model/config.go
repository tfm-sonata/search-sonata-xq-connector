package model

var RequestConfig *Config

type Config struct {
	// channel mame
	Id                       string `json:"id"`
	Info                     string `json:"INFO"`
	SqsBaseUrl               string `json:"SQS_BASE_URL"`
	QueueNameOut             string `json:"QUEUE_OUT"`
	QueueNameTracking        string `json:"QUEUE_NAME_TRACKING"`
	WebserviceUrl            string `json:"WS_URL"`
	WebserviceUser           string `json:"WS_USER"`
	WebservicePassword       string `json:"WS_PASSWORD"`
	WebserviceLevel1         string `json:"WS_LEVEL1"`
	WebserviceSessionTimeout string `json:"WS_SESSION_TIMEOUT"`
	WebserviceRequestType    string `json:"WS_REQUESTTYPE"`
	WebserviceAgent          string `json:"WS_AGENT"`
	WebserviceInstance       string `json:"WS_INSTANCE"`

	//headers
	WebserviceParty            string `json:"WS_PARTY"`
	WebserviceConsumer         string `json:"WS_CONSUMER"`
	WebserviceConsumerLocation string `json:"WS_CONSUMER_LOCATION"`
	WebserviceConsumerType     string `json:"WS_CONSUMER_TYPE"`
	WebserviceConsumerTime     string `json:"WS_CONSUMER_TIME"`
	WebserviceRelationshipType string `json:"WS_RELATIONSHIPTYPE"`
	WebserviceRelatesToValue   string `json:"WS_RELATESTOVALUE"`
	WebserviceMessageId        string `json:"WS_MESSAGE_ID"`

	//AirShppoingRequest
	WebserviceVersion         string `json:"WS_VERSION"`
	WebserviceCorrelationID   string `json:"WS_CORRELATION_ID"`
	WebserviceProperName      string `json:"WS_PROPER_NAME"`
	WebservicePseudoCity      string `json:"WS_PSEUDOCITY"`
	WebserviceIATA_Nbr        string `json:"WS_IATA_NBR"`
	WebserviceSequence        int    `json:"WS_SEQUENCE"`
	WebserviceName            string `json:"WS_NAME"`
	WebserviceSystemID        string `json:"WS_SYSTEM_ID"`
	WebserviceAirLineID       string `json:"WS_AIRLINE_ID"`
	WebserviceAirLineName     string `json:"WS_AIRLINE_NAME"`
	WebserviceParticipantName string `json:"WS_PARTICIPANT_NAME"`

	WebserviceAPIKey                string `json:"WS_API_KEY"`
	WebserviceSOAPAction            string `json:"WS_SOAP_ACTION"`
	WebserviceFarePreferencesCotext string `json:"WS_FAREPREFERENCES_COTEXT"`
	WebserviceFarePreferencesCode   string `json:"WS_FAREPREFERENCES_CODE"`
	WebserviceFareCode              string `json:"WS_FARE_CODE"`
	WebserviceCabinCode             string `json:"WS_CABIN_CODE"`
}
