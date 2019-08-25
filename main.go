package main

import (
	"os"
	"strings"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func RetrieveSecret(variableName string) {
	// All clients require a Session. The Session provides the client with
	// shared configuration such as region, endpoint, and credentials. A
	// Session should be shared where possible to take advantage of
	// configuration and credential caching. See the session package for
	// more information.
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		printAndExit(err)
	}

	// AWS Go SDK does not currently support automatic fetching of region from ec2metadata.
	// If the region could not be found in an environment variable or a shared config file,
	// create metaSession to fetch the ec2 instance region and pass to the regular Session.
	if *sess.Config.Region == "" {
		metaSession, err := session.NewSession()
		if err != nil {
			printAndExit(err)
		}

		metaClient := ec2metadata.New(metaSession)
		// If running on an EC2 instance, the metaClient will be available and we can set the region to match the instance
		// If not on an EC2 instance, the region will remain blank and AWS returns a "MissingRegion: ..." error
		if metaClient.Available() {
			if region, err := metaClient.Region(); err == nil {
				sess.Config.Region = aws.String(region)
			} else {
				printAndExit(err)
			}
		}
	}

	// Create a new instance of the service's client with a Session.
	// Optional aws.Config values can also be provided as variadic arguments
	// to the New function. This option allows you to provide service
	// specific configuration.

	svc := ssm.New(sess)

	// Check if key has been specified
	arguments := strings.SplitN(variableName, "#", 2)

	secretName := arguments[0]

	// Get secret value
	withDecryption := true
	param, err := svc.GetParameter(&ssm.GetParameterInput{
		Name: aws.String(secretName),
		WithDecryption: &withDecryption,
	})

	if err != nil {
		printAndExit(err)
	}

	fmt.Print(string(*param.Parameter.Value))

}

func main() {
	if len(os.Args) != 2 {
		os.Stderr.Write([]byte("A variable ID or version flag must be given as the first and only argument!\n"))
		os.Exit(-1)
	}

	// Get the secret and key name from the argument
	singleArgmument := os.Args[1]

	switch singleArgmument {
	case "-v", "--version":
		os.Stdout.Write([]byte(VERSION))
	default:
		RetrieveSecret(singleArgmument)
	}
}

func printAndExit(err error) {
	os.Stderr.Write([]byte(err.Error()))
	os.Exit(1)
}
