// Package soap provides a SOAP HTTP client.
package wsdl2goEdit

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
)

// XSINamespace is a link to the XML Schema instance namespace.
const XSINamespace = "http://www.w3.org/2001/XMLSchema-instance"

var xmlTyperType reflect.Type = reflect.TypeOf((*XMLTyper)(nil)).Elem()

// A RoundTripper executes a request passing the given req as the SOAP
// envelope body. The HTTP response is then de-serialized onto the resp
// object. Returns error in case an error occurs serializing req, making
// the HTTP request, or de-serializing the response.
type RoundTripper interface {
	RoundTrip(req, resp Message) error
	RoundTripSoap12(action string, req, resp Message) error
}

// Message is an opaque type used by the RoundTripper to carry XML
// documents for SOAP.
type Message interface{}

// Header is an opaque type used as the SOAP Header element in requests.
type Header interface{}

// AuthHeader is a Header to be encoded as the SOAP Header element in
// requests, to convey credentials for authentication.
type AuthHeader struct {
	Namespace string `xml:"xmlns:ns,attr"`
	Username  string `xml:"ns:username"`
	Password  string `xml:"ns:password"`
}

// Client is a SOAP client.
type Client struct {
	URL                    string               // URL of the server
	UserAgent              string               // User-Agent header will be added to each request
	Namespace              string               // SOAP Namespace
	URNamespace            string               // Uniform Resource Namespace
	ThisNamespace          string               // SOAP This-Namespace (tns)
	DefaultNamespace       string               // The Default
	ExcludeActionNamespace bool                 // Include Namespace to SOAP Action header
	Envelope               string               // Optional SOAP Envelope
	Header                 Header               // Optional SOAP Header
	ContentType            string               // Optional Content-Type (default text/xml)
	Config                 *http.Client         // Optional HTTP client
	Pre                    func(*http.Request)  // Optional hook to modify outbound requests
	Post                   func(*http.Response) // Optional hook to snoop inbound responses
}

// XMLTyper is an abstract interface for types that can set an XML type.
type XMLTyper interface {
	SetXMLType()
}

func setXMLType(v reflect.Value) {
	if !v.IsValid() {
		return
	}
	switch v.Type().Kind() {
	case reflect.Interface:
		setXMLType(v.Elem())
	case reflect.Ptr:
		if v.IsNil() {
			break
		}
		ok := v.Type().Implements(xmlTyperType)
		if ok {
			v.MethodByName("SetXMLType").Call(nil)
		}
		setXMLType(v.Elem())
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			setXMLType(v.Index(i))
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).CanAddr() {
				setXMLType(v.Field(i).Addr())
			} else {
				setXMLType(v.Field(i))
			}
		}
	}
}

func doRoundTrip(c *Client, setHeaders func(*http.Request), in, out Message) error {
	setXMLType(reflect.ValueOf(in))
	req := &Envelope{
		EnvelopeAttr: c.Envelope,
		URNAttr:      c.URNamespace,
		NSAttr:       c.Namespace,
		TNSAttr:      c.ThisNamespace,
		XSIAttr:      XSINamespace,
		Header:       c.Header,
		Body:         in,
		Default:      c.DefaultNamespace,
	}

	if req.EnvelopeAttr == "" {
		req.EnvelopeAttr = "http://schemas.xmlsoap.org/soap/envelope/"
	}
	if req.NSAttr == "" {
		req.NSAttr = c.URL
	}
	if req.TNSAttr == "" {
		req.TNSAttr = req.NSAttr
	}
	var b bytes.Buffer
	err := xml.NewEncoder(&b).Encode(req)
	if err != nil {
		return err
	}
	cli := c.Config
	if cli == nil {
		cli = http.DefaultClient
	}
	r, err := http.NewRequest("POST", c.URL, &b)
	if err != nil {
		return err
	}
	setHeaders(r)
	if c.Pre != nil {
		c.Pre(r)
	}
	printWsdl(req)
	resp, err := cli.Do(r)
	if err != nil {
		return err
	}
	//writeResponseWsdl(resp.Body)
	defer resp.Body.Close()
	if c.Post != nil {
		c.Post(resp)
	}
	if resp.StatusCode != http.StatusOK {
		// read only the first MiB of the body in error case
		limReader := io.LimitReader(resp.Body, 1024*1024)
		body, _ := ioutil.ReadAll(limReader)
		return &HTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Msg:        string(body),
		}
	}

	marshalStructure := struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    Message
	}{Body: out}
	return xml.NewDecoder(resp.Body).Decode(&marshalStructure)
}

// RoundTrip implements the RoundTripper interface.
func (c *Client) RoundTrip(in, out Message) error {
	headerFunc := func(r *http.Request) {
		if c.UserAgent != "" {
			r.Header.Add("User-Agent", c.UserAgent)
		}
		/*var actionName, soapAction string
		if in != nil {
			soapAction = reflect.TypeOf(in).Elem().Name()
		}*/
		ct := c.ContentType
		if ct == "" {
			ct = "text/xml"
		}
		r.Header.Set("Content-Type", ct)
		/*if in != nil {
			if c.ExcludeActionNamespace {
				actionName = soapAction
			} else {
				actionName = fmt.Sprintf("%s/%s", c.Namespace, soapAction)
			}
			r.Header.Add("SOAPAction", actionName)
		}*/
	}
	return doRoundTrip(c, headerFunc, in, out)
}

// RoundTripWithAction implements the RoundTripper interface for SOAP clients
// that need to set the SOAPAction header.
func (c *Client) RoundTripWithAction(soapAction string, in, out Message) error {
	headerFunc := func(r *http.Request) {
		if c.UserAgent != "" {
			r.Header.Add("User-Agent", c.UserAgent)
		}
		//var actionName string
		ct := c.ContentType
		if ct == "" {
			ct = "text/xml"
		}
		r.Header.Set("Content-Type", ct)
		/*if in != nil {
			if c.ExcludeActionNamespace {
				actionName = soapAction
			} else {
				actionName = fmt.Sprintf("%s/%s", c.Namespace, soapAction)
			}
			r.Header.Add("SOAPAction", actionName)
		}*/
	}
	return doRoundTrip(c, headerFunc, in, out)
}

// RoundTripSoap12 implements the RoundTripper interface for SOAP 1.2.
func (c *Client) RoundTripSoap12(action string, in, out Message) error {
	headerFunc := func(r *http.Request) {
		r.Header.Add("Content-Type", fmt.Sprintf("application/soap+xml; charset=utf-8; action=\"%s\"", action))
	}
	return doRoundTrip(c, headerFunc, in, out)
}

// HTTPError is detailed soap http error
type HTTPError struct {
	StatusCode int
	Status     string
	Msg        string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("%q: %q", e.Status, e.Msg)
}

// Envelope is a SOAP envelope.
type Envelope struct {
	XMLName      xml.Name `xml:"SOAP-ENV:Envelope"`
	EnvelopeAttr string   `xml:"xmlns:SOAP-ENV,attr"`
	NSAttr       string   `xml:"xmlns:ns,attr"`
	TNSAttr      string   `xml:"xmlns:tns,attr,omitempty"`
	URNAttr      string   `xml:"xmlns:urn,attr,omitempty"`
	XSIAttr      string   `xml:"xmlns:xsi,attr,omitempty"`
	Header       Message  `xml:"SOAP-ENV:Header"`
	Body         Message  `xml:"SOAP-ENV:Body"`
	Default      string   `xml:"xmlns,attr"`
}

//this funtion is used for testing request xml generation
func printWsdl(request interface{}) {
	var w io.Writer
	path, err := filepath.Abs("../resources/test/request.xml")
	if err != nil {
		log.Print("something went wrong \n", err)
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w = f
	if err != nil {
		//log.Errorf("search(%+v) failed.", err)
	}

	enc := xml.NewEncoder(w)
	enc.Indent("  ", "    ")
	if err := enc.Encode(request); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

//this funtion is used for testing response xml generation
func writeResponseWsdl(src io.ReadCloser) {
	path, err := filepath.Abs("../resources/test/response.xml")
	if err != nil {
		log.Print("something went wrong \n", err)
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Print("printing to wsdl error")
		log.Fatal(err)
	}
	defer f.Close()

	io.Copy(f, src)
}
