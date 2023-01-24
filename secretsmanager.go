package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"log"
	"time"
)

func loadSecretBinary(secretsClient *secretsmanager.Client, secretName *string, timeoutSeconds int) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	req := &secretsmanager.GetSecretValueInput{
		SecretId: secretName,
	}

	result, err := secretsClient.GetSecretValue(ctx, req)
	if err != nil {
		return nil, err
	} else {
		log.Printf("INFO Obtained version %d of secret %s modified at %s", result.VersionId, *result.Name, result.CreatedDate.Format(time.RFC3339))
		return result.SecretBinary, nil
	}
}

func loadSecretString(secretsClient *secretsmanager.Client, secretName *string, timeoutSeconds int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	req := &secretsmanager.GetSecretValueInput{
		SecretId: secretName,
	}

	result, err := secretsClient.GetSecretValue(ctx, req)
	if err != nil {
		return "", err
	} else {
		log.Printf("INFO Obtained version %d of secret %s modified at %s", result.VersionId, *result.Name, result.CreatedDate.Format(time.RFC3339))
		return *result.SecretString, nil
	}
}

func loadSecretB64(secretsClient *secretsmanager.Client, secretName *string, timeoutSeconds int) (string, error) {
	bytes, err := loadSecretBinary(secretsClient, secretName, timeoutSeconds)
	if err != nil {
		return "", err
	} else {
		return base64.StdEncoding.EncodeToString(bytes), nil
	}
}

func loadGoogleAuthSecretJson(secretsClient *secretsmanager.Client, secretName *string, timeoutSeconds int) (*GoogleAuthSecretJson, error) {
	bytes, err := loadSecretString(secretsClient, secretName, timeoutSeconds)
	if err != nil {
		return nil, err
	}

	var result GoogleAuthSecretJson
	err = json.Unmarshal([]byte(bytes), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
