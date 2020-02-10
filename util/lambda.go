package util

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/tespo/satya/v2/types"
)

//
// TriggerLambda will trigger a lambda function
//
func TriggerLambda(payload types.Payload, name string) (*lambda.InvokeOutput, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := lambda.New(sess, &aws.Config{Region: aws.String("us-east-1")})

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	result, err := client.Invoke(&lambda.InvokeInput{FunctionName: aws.String(name), Payload: data})
	if err != nil {
		return nil, err
	}

	return result, nil
}
