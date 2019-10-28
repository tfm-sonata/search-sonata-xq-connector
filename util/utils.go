package util

import (
	"encoding/json"
	"fmt"
	"math"
)

// a function print well formatted jsons
func Round(x, unit float64) float64 {
	return math.Round(x*unit) / unit
}

func PrettyPrint(jsonBody interface{}) {
	prettyData, _ := json.MarshalIndent(&jsonBody, "", "  ")
	fmt.Printf("%s\n", prettyData)
}

// a returns a well formatted json
func MapToString(jsonBody interface{}) string {
	s, _ := json.Marshal(&jsonBody)
	return fmt.Sprintf("%s", s)
}

// a returns a well formatted json
func MapToPrettyString(jsonBody interface{}) string {
	prettyData, _ := json.MarshalIndent(&jsonBody, "", "  ")
	return fmt.Sprintf("%s", prettyData)
}
func MarshalJson(jsonBody interface{}) (string, error) {
	s, err := json.Marshal(&jsonBody)
	return fmt.Sprintf("%s", s), err
}
