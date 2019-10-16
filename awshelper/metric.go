package awshelper

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"log"
	"strings"
)

const (
	LAMBDA            = "LAMBDA"
	SERVICE           = "SERVICE"
	THIRDPARTYSERVICE = "THIRDPARTYSERVICE"
	RECEIVED_MESSAGES = "RECEIVED_MESSAGES"
	RESPONSE_TIME     = "RESPONSE_TIME"
	UNIT_COUNT        = "Count"
	UNIT_MILLISECONDS = "Milliseconds"
)

var (
	cwClient *cloudwatch.CloudWatch
)

// AddMetric adds a service metrics to cloudwatch Metric
func AddMetric(metricName, unit string, value float64) {
	addMetric(SERVICE, ServiceName, metricName, unit, value)
}

// AddThirdPartyMetric adds the 3rd party service metrics to cloudwatch Metric
func AddThirdPartyMetric(metricName, thirdPartyServiceName, unit string, value float64) {
	addMetric(THIRDPARTYSERVICE, thirdPartyServiceName, metricName, unit, value)
}

func addMetric(dimensionName, dimensionValue, metricName string, unit string, value float64) {

	// if required create new client to handle cloudwatch events
	if cwClient == nil {
		cwClient = cloudwatch.New(AWSSession())
	}

	_, err := cwClient.PutMetricData(&cloudwatch.PutMetricDataInput{
		Namespace: aws.String(LAMBDA),
		MetricData: []*cloudwatch.MetricDatum{
			&cloudwatch.MetricDatum{
				MetricName: aws.String(strings.ToUpper(metricName)),
				Unit:       aws.String(unit),
				Value:      aws.Float64(value),
				Dimensions: []*cloudwatch.Dimension{
					&cloudwatch.Dimension{
						Name:  aws.String(dimensionName),
						Value: aws.String(dimensionValue),
					},
				},
			},
		},
	})
	if err != nil {
		log.Println("Couldn't put metric:", err)
		return
	}
	log.Println(dimensionName+" metric:", metricName, "successfull added. Value", value, unit)
}
