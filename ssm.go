package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"log"
	"regexp"
	"time"
)

func GetSSMValue(ssmClient *ssm.Client, app *AppId, subkey string, timeoutSeconds int) (*types.Parameter, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	defer cancelFunc()

	req := &ssm.GetParameterInput{
		Name:           aws.String(app.SsmPath(subkey)),
		WithDecryption: aws.Bool(true),
	}
	response, err := ssmClient.GetParameter(ctx, req)
	if err != nil {
		return nil, err
	}
	log.Printf("INFO Got version %d of parameter %s last modified at %s", response.Parameter.Version, *response.Parameter.Name, response.Parameter.LastModifiedDate.Format(time.RFC3339))
	return response.Parameter, nil
}

func GetSSMValueString(ssmClient *ssm.Client, app *AppId, subkey string, timeoutSeconds int) (string, error) {
	param, err := GetSSMValue(ssmClient, app, subkey, timeoutSeconds)
	if err != nil {
		return "", err
	} else {
		if param.Type != types.ParameterTypeString {
			log.Printf("WARNING Parameter %s is not a string, it's type %s", *param.Name, param.Type)
		}
		return *param.Value, nil
	}
}

func GetSSMValueStringList(ssmClient *ssm.Client, app *AppId, subkey string, timeoutSeconds int) ([]string, error) {
	param, err := GetSSMValue(ssmClient, app, subkey, timeoutSeconds)
	if err != nil {
		return nil, err
	} else {
		if param.Type != types.ParameterTypeStringList {
			log.Printf("WARNING Parameter %s is not a string list, it's type %s", *param.Name, param.Type)
		}
		splitter := regexp.MustCompile("\\s*,\\s*")

		return splitter.Split(*param.Value, -1), nil
	}
}
