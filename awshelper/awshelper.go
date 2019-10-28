package awshelper

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"log"
	"os"
)

var (
	ServiceName = os.Getenv("SERVICENAME")
	awsSession  *session.Session
)

func AWSSession() *session.Session {
	if awsSession == nil {
		CreateSession()
	}
	var seesionDetails credentials.Value
	seesionDetails, _ = awsSession.Config.Credentials.Get()
	log.Println("session created ", seesionDetails)
	return awsSession
}

// creating a session for all aws services
func CreateSession() {

	const MT = "init"
	// Initialize a session that the SDK uses to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and configuration from the shared configuration file ~/.aws/config.
	/*s, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           "tfm-develop",
	})*/

	s, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)

	if err != nil {
		log.Println("Couldn't create AWS session", err)
	}
	awsSession = s
	log.Println("AWS Session successfully created.")

}

var lambdaClient lambdaiface.LambdaAPI

func GetLambdaClient() lambdaiface.LambdaAPI {
	// the
	if lambdaClient == nil {
		lambdaClient = lambda.New(AWSSession())
	}

	return lambdaClient
}

func LambdaClient(lc lambdaiface.LambdaAPI) {
	lambdaClient = lc
}
