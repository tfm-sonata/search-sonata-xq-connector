package search

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"
	"time"
)

type Notes struct {
	To      string `xml:"to"`
	From    string `xml:"from"`
	Heading string `xml:"heading"`
	Body    string `xml:"body"`
}

type Message interface{}

func TestUnmarshalAirShoppingRS(t *testing.T) {
	//t.Skip()
	path, err := filepath.Abs("../resources/test/response.xml")
	if err != nil {
		log.Print("something went wrong \n", err)
	}
	data, _ := ioutil.ReadFile(path)

	/*γ := struct {
		AirShoppingResponse `xml:"AirShoppingResponse"`
	}{}*/
	γ := struct {
		AirShoppingResponse `xml:"AirShoppingResponse"`
	}{}
	marshalStructure := struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    Message
	}{Body: &γ}

	error := xml.Unmarshal([]byte(data), &marshalStructure)

	assert.Nil(t, error, "unmarshall failed \n", error)

	if nil != error {
		fmt.Println("Error in unmarshalling the response ", error)
	}
	response := &marshalStructure.Body
	if response != nil {
		fmt.Println("Response ", response)
	}

	//log.Printf(" No. of pax  %+v\n", γ.AirShoppingResponse.AirShoppingRS.OffersGroup.AllOffersSnapshot.PassengerQuantity)
	tfmMapper := NewTfmMapper()
	tfmResponse, err := tfmMapper.CreateTFMResponse(γ.AirShoppingRS, "", "", time.Duration(10), "")
	log.Println("Tested mapping !!")

	searchResponse, error := json.Marshal(tfmResponse)
	if error != nil {
		log.Println("unable to marsha Tfm response", error)
	}

	fmt.Println(string(searchResponse))

}
