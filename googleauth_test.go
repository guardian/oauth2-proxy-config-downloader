package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestGoogleAuthGenerate(t *testing.T) {
	expectedContent := `{"type":"service_account","project_id":"some-project","private_key_id":"some-key-id","private_key":"PRIVATEKEYVALUE","client_email":"clientemail@something","client_id":"some-client-id","auth_uri":"https://some-auth-uri","token_uri":"https://some-token-uri","auth_provider_x509_cert_url":"https://cert-provider-uri","client_x509_cert_url":"https://client-cert-uri"}`

	fakeGoogleConfig := &GoogleAuthSecretJson{
		Type:                    "service_account",
		ProjectId:               "some-project",
		PrivateKeyId:            "some-key-id",
		PrivateKey:              "PRIVATEKEYVALUE",
		ClientEmail:             "clientemail@something",
		ClientId:                "some-client-id",
		AuthUri:                 "https://some-auth-uri",
		TokenUri:                "https://some-token-uri",
		AuthProviderX509CertUrl: "https://cert-provider-uri",
		ClientX509CertUrl:       "https://client-cert-uri",
	}

	err := fakeGoogleConfig.Output("/tmp")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fp, err := os.Open("/tmp/auth.json")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer fp.Close()
	defer os.Remove("/tmp/auth.json")

	content, err := ioutil.ReadAll(fp)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if string(content) != expectedContent {
		t.Error("Rendered content did not match expected")
		t.Error("Expected: ", expectedContent)
		t.Error("Got: ", string(content))
	}
}
