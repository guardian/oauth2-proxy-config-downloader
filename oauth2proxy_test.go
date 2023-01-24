package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestOauth2ProxyGenerate(t *testing.T) {
	expectedResult := `#################################
#  Auto-generated configuration for OAuth2Proxy
#################################
OPTS="      --upstream=http://localhost:9000 \
      --upstream=file:///var/www/healthcheck/#/healthcheck/ \
      --skip-auth-regex=/healthcheck/_ \
      --redirect-url=https://app.my-domain.co.uk/oauth2/callback \
      --http-address=0.0.0.0:4180 \
      --google-service-account-json=/tmp/auth.json \
      --google-group=group-one \
      --google-admin-email=google@admin.address \
      --email-domain=my-domain.co.uk \
      --cookie-secret=generated-secret \
      --cookie-name=my-cookie-name \
      --client-secret=someClientSecret \
      --client-id=someClientId"
`

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

	config := NewOAuth2ProxyConfig(fakeGoogleConfig,
		"someClientId",
		"someClientSecret",
		"google@admin.address",
		[]string{"group-one"},
		"my-domain.co.uk",
		"generated-secret",
		"my-cookie-name",
		"app.my-domain.co.uk",
		[]string{"http://localhost:9000"},
		"/tmp",
	)

	err := config.Generate("/tmp")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fp, err := os.Open("/tmp/oauth2proxy.env")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer fp.Close()
	defer os.Remove("/tmp/auth2proxy.env")

	content, err := ioutil.ReadAll(fp)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	stringContent := string(content)
	if stringContent != expectedResult {
		t.Error("Generated env did not match expected result")
		t.Error("Got: ", stringContent)
		t.Error("Expected result was: ", expectedResult)
	}
}
