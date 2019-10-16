package awshelper

import (
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/search-sonata-xq-connector/util"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"log"
)

var lambdaServiceClient lambdaiface.LambdaAPI

func InvokeLambdaResponseFunction(functionName string, payload []byte) ([]byte, error) {
	return invokeLambdaFunction(functionName, payload, lambda.InvocationTypeRequestResponse)
}

func InvokeLambdaEventFunction(functionName string, payload []byte) ([]byte, error) {
	return invokeLambdaFunction(functionName, payload, lambda.InvocationTypeEvent)
}

// this InvokeLambdaResponseFunction function handles all LambdaInvocation Requests
func invokeLambdaFunction(functionName string, payload []byte, iType string) ([]byte, error) {

	// create an Lambda InvokeInput
	input := &lambda.InvokeInput{
		FunctionName:   &functionName,
		InvocationType: &iType,
		Payload:        payload,
	}

	lambdaServiceClient = GetLambdaClient()
	log.Print("invoking lambda function")
	output, err := lambdaServiceClient.Invoke(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case lambda.ErrCodeServiceException:
				fmt.Println(lambda.ErrCodeServiceException, aerr.Error())
			case lambda.ErrCodeResourceNotFoundException:
				fmt.Println(lambda.ErrCodeResourceNotFoundException, aerr.Error())
			case lambda.ErrCodeInvalidRequestContentException:
				fmt.Println(lambda.ErrCodeInvalidRequestContentException, aerr.Error())
			case lambda.ErrCodeInvalidRuntimeException:
				fmt.Println(lambda.ErrCodeInvalidRuntimeException, aerr.Error())
			default:
				fmt.Println(" default error ", aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Printf("err.Error() @ %s:  %s\n", *input.FunctionName, err.Error())
		}

		//	_ = SendTrackingMessage(TrackingMessage{Title: "Error Invoking Lambda: " + *input.FunctionName, Action: "ERROR", ActionType: "RESPONSE", ContentType: "JSON", BodyZ: util.MapToPrettyString(&err), ServiceName: "BOOKING-FRAUDCHECK"})
		return nil, err
	}

	fmt.Printf("Result from LambdaInvokationRequest %s, returned:%s\n", *input.FunctionName, util.MapToString(output.Payload))
	return output.Payload, nil
}
