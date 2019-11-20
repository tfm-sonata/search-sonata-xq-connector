package search

import (
	"crypto/tls"
	"errors"
	"git-codecommit.eu-central-1.amazonaws.com/search-sonata-xq-connector/model"
	"git-codecommit.eu-central-1.amazonaws.com/search-sonata-xq-connector/util"
	"git-codecommit.eu-central-1.amazonaws.com/search-sonata-xq-connector/wsdl2goEdit"
	"log"
	"net/http"
	"time"
)

func init() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	http.DefaultClient.Timeout = 17 * time.Second
}

type SoapHandler interface {
	Search(query Query) (*TfmResponse, error)
}

type SoapHandlerImpl struct {
	config    WebserviceConfig
	tfmMapper TfmMapper
	helper    SoapHandlerHelper
}

func NewSoapHandler(cfg WebserviceConfig) SoapHandler {
	tfmMapper := NewTfmMapper()
	return &SoapHandlerImpl{
		config:    cfg,
		tfmMapper: tfmMapper,
		helper:    NewSoapHandlerHelper(),
	}
}

var (
	q      Query
	cookie string
)

type WebserviceConfig struct {
	WsUrl            string
	WsUser           string
	WsPassword       string
	Level1           string
	WsSessionTimeout string
	ReqModus         string
	Agent            string
	Instance         string
}

func (this *SoapHandlerImpl) Search(query Query) (*TfmResponse, error) {
	q = query
	cookie = ""
	helper := this.helper
	searchFactory := helper.SearchFactory()

	securityHeader, _ := searchFactory.createSOAPHeader(model.RequestConfig)
	airShoppingRQ, _ := searchFactory.CreateAirShoppingRQ(query)

	soapService := helper.createProviderService(this.config, securityHeader)

	log.Println("Executing ProvideAirShopping..")
	timeStart := time.Now()

	var airShoppingRS *AirShoppingRS
	var err error

	airShoppingRS, err = soapService.ShopAir(airShoppingRQ)
	if err != nil {
		log.Println("Something went wrong during ProvideAirShopping. Possible Timeout. Returning empty tfmResponse.", err)
		return nil, err
	}

	err = helper.validateErrors(airShoppingRS)
	if err != nil {
		log.Println("Validation error. Returning empty tfmResponse.", err)
		return nil, nil
	}
	//var correlationId = string(*airShoppingRQ.CorrelationID)
	timeEnd := time.Now()
	// defer awshelper.AddThirdPartyMetric(awshelper.RESPONSE_TIME, "Mercado FlightVacancyNpmRQ", awshelper.UNIT_MILLISECONDS, float64(timeEnd.Sub(timeStart)/time.Millisecond))
	log.Println("ProvideAirShopping finished")

	tfmResponse, err := this.tfmMapper.CreateTFMResponse(airShoppingRS, "", cookie, timeEnd.Sub(timeStart)/time.Millisecond)

	return tfmResponse, err
}
func interceptRequest(req *http.Request) {
	//req.Header.Set("Accept", "*/*")
	req.Header.Set("Username", model.RequestConfig.WebserviceUser)
	req.Header.Set("Content-Type", model.RequestConfig.WebServiceContentType)
	req.Header.Set("Password", model.RequestConfig.WebservicePassword)
	//req.Header.Set("SOAPAction", strings.Join([]string{"`\"", model.RequestConfig.WebserviceSOAPAction, "\"`"}, ""))
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	/*soapAction := req.Header.Get("Soapaction")
	url := req.URL

	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Println("Couldn't dump httpRequest for tracking", err)
	}
	awshelper.SendTrackingMessage(awshelper.TrackingMessage{Timestamp: aws.TimeUnixMilli(time.Now()), ServiceName: awshelper.ServiceName, Action: "API-Call", SessionId: q.SessionId, RequestId: q.RequestId, ContentType: "XML", BodyZ: string(requestDump), Title: soapAction + ", URL: " + url.Host, ActionType: "Request"})*/

}

func interceptResponse(resp *http.Response) {
	if cookie == "" {
		cookie = resp.Header.Get("Set-Cookie")
	}
	/*responseDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Println("Couldn't dump httpResponse for tracking", err)
	}
	awshelper.SendTrackingMessage(awshelper.TrackingMessage{Timestamp: aws.TimeUnixMilli(time.Now()), ServiceName: awshelper.ServiceName, Action: "API-Call", SessionId: q.SessionId, RequestId: q.RequestId, ContentType: "XML", BodyZ: string(responseDump), Title: "Webservice Response", ActionType: "Response"})*/
}

func createClient(config WebserviceConfig, header wsdl2goEdit.Header) wsdl2goEdit.Client {
	log.Println("Creating client..")
	//proxyURL, _ := url.Parse("http://10.145.10.5:8080")
	cli := wsdl2goEdit.Client{
		URL:                    config.WsUrl,
		Header:                 header,
		DefaultNamespace:       "http://interes.com/clients/tui/webservice/types/v140",
		Namespace:              "http://interes.com/clients/tui/queueschubser/service/v140/types",
		ContentType:            "text/xml",
		ExcludeActionNamespace: true,
		Pre:                    interceptRequest,
		Post:                   interceptResponse,
		//TODO: this config needs to be reomoved in final code . this works only in sonata network
		/*Config: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
				Proxy: http.ProxyURL(proxyURL),
			},
		},*/
	}
	return cli
}

func NewSoapHandlerHelper() SoapHandlerHelper {
	return &SoapHandlerHelperImpl{}
}

type SoapHandlerHelper interface {
	SearchFactory() SearchFactory
	validateErrors(airShoppingRS *AirShoppingRS) error
	createProviderService(config WebserviceConfig, header interface{}) NdcResService17_2Port
}
type SoapHandlerHelperImpl struct {
}

func (*SoapHandlerHelperImpl) SearchFactory() SearchFactory {
	return &SearchFactoryImpl{}
}

func (*SoapHandlerHelperImpl) validateErrors(airShoppingRS *AirShoppingRS) error {
	if airShoppingRS.Errors != nil {
		out, _ := util.MarshalJson(airShoppingRS.Errors.Error.Value)
		return errors.New(string(out))
	}
	return nil
}
func (this *SoapHandlerHelperImpl) createProviderService(config WebserviceConfig, header interface{}) NdcResService17_2Port {
	client := createClient(config, header)
	return NewNdcResService17_2Port(&client)
}
