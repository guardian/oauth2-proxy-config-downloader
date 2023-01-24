package main

import (
	"encoding/json"
	"os"
	"path"
)

type GoogleWebClient struct {
	ClientId                string   `json:"client_id"`
	ProjectId               string   `json:"project_id"`
	AuthUri                 string   `json:"auth_uri"`
	TokenUri                string   `json:"token_uri"`
	AuthProviderX509CertUrl string   `json:"auth_provider_x509_cert_url"`
	ClientSecret            string   `json:"client_secret"`
	RedirectUris            []string `json:"redirect_uris"`
	JavascriptOrigins       []string `json:"javascript_origins"`
}

type GoogleAuthSecretJson struct {
	Type                    string `json:"type"`
	ProjectId               string `json:"project_id"`
	PrivateKeyId            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientId                string `json:"client_id"`
	AuthUri                 string `json:"auth_uri"`
	TokenUri                string `json:"token_uri"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl       string `json:"client_x509_cert_url"`
}

/*
Validate returns true if all the required fields in the json are populated,
or false otherwise.
TODO: this validation is rather rough at the moment!
*/
func (j *GoogleAuthSecretJson) Validate() bool {
	return j.ClientId != "" &&
		j.ProjectId != "" &&
		j.AuthUri != "" &&
		j.TokenUri != "" &&
		j.AuthProviderX509CertUrl != "" &&
		j.PrivateKeyId != "" &&
		j.Type != "" &&
		j.ClientEmail != ""
}

func (j *GoogleAuthSecretJson) Output(dir string) error {
	outputFilename := path.Join(dir, "auth.json")
	contentToWrite, _ := json.Marshal(j) //no way that this should error as it's all simple types and has been parsed in
	fp, err := os.OpenFile(outputFilename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer fp.Close()

	_, err = fp.Write(contentToWrite)
	return err
}
