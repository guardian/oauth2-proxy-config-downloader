package main

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"regexp"
	"sort"
)

type OAuth2ProxyConfig struct {
	GoogleAdminEmail         string   `json:"google-admin-email"`
	GoogleServiceAccountJson string   `json:"google-service-account-json"` //path to the JSON file, that is
	AllowedGoogleGroups      []string `json:"google-group"`
	EmailDomain              string   `json:"email-domain"`
	ClientID                 string   `json:"client-id"`
	ClientSecret             string   `json:"client-secret"`
	CookieName               string   `json:"cookie-name"`
	CookieSecret             string   `json:"cookie-secret"`
	HttpAddress              string   `json:"http-address"`
	RedirectUri              string   `json:"redirect-url"`
	SkipAuthRegex            string   `json:"skip-auth-regex"`
	Upstreams                []string `json:"upstream"`
}

/*
This function will help you to convert your object from struct to map[string]interface{} based on your JSON tag in your structs.
Thanks to https://gist.github.com/bxcodec/c2a25cfc75f6b21a0492951706bc80b8
*/
func structToMap(item interface{}) map[string]interface{} {

	res := map[string]interface{}{}
	if item == nil {
		return res
	}
	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get("json")
		field := reflectValue.Field(i).Interface()
		if tag != "" && tag != "-" {
			if v.Field(i).Type.Kind() == reflect.Struct {
				res[tag] = structToMap(field)
			} else {
				res[tag] = field
			}
		}
	}
	return res
}

func NewOAuth2ProxyConfig(googleServiceAccountJson *GoogleAuthSecretJson, googleClientId string, googleClientSecret string, googleAdminEmail string, googleGroups []string, emailDomain string, cookieSecret string, cookieName string, domainName string, upstreamUriList []string, outputDir string) *OAuth2ProxyConfig {
	truncatedCookieSecret := cookieSecret
	if len(cookieSecret) > 32 {
		truncatedCookieSecret = cookieSecret[0:32]
	}

	actualGoogleClientId := googleServiceAccountJson.ClientId
	if googleClientId != "" {
		actualGoogleClientId = googleClientId
	}
	return &OAuth2ProxyConfig{
		GoogleAdminEmail:         googleAdminEmail,
		GoogleServiceAccountJson: path.Join(outputDir, "auth.json"),
		AllowedGoogleGroups:      googleGroups,
		EmailDomain:              emailDomain,
		ClientID:                 actualGoogleClientId,
		ClientSecret:             googleClientSecret,
		CookieSecret:             truncatedCookieSecret,
		CookieName:               cookieName,
		HttpAddress:              "0.0.0.0:4180",
		RedirectUri:              fmt.Sprintf("https://%s/oauth2/callback", domainName),
		SkipAuthRegex:            "/healthcheck/_",
		Upstreams:                append([]string{"file:///var/www/healthcheck/#/healthcheck/"}, upstreamUriList...),
	}
}

/*
Generate outputs an oauth2proxy environment file to the given directory
*/
func (p *OAuth2ProxyConfig) Generate(outputPath string) error {
	outputFile := path.Join(outputPath, "oauth2proxy.env")

	rows := make([]string, 0)

	mapData := structToMap(p)
	for k, v := range mapData {
		typeof := reflect.TypeOf(v)
		if typeof.Kind() == reflect.Slice {
			array := reflect.ValueOf(v)
			for i := 0; i < array.Len(); i++ {
				rows = append(rows, fmt.Sprintf("--%s=%s", k, array.Index(i)))
			}
		} else {
			rows = append(rows, fmt.Sprintf("--%s=%s", k, v))
		}
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i] > rows[j]
	})

	fp, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer fp.Close()

	isEmpty := regexp.MustCompile("^\\s*$")

	fp.WriteString("#################################\n")
	fp.WriteString("#  Auto-generated configuration for OAuth2Proxy\n")
	fp.WriteString("#################################\n")
	fp.WriteString("OPTS=\"")
	rowCount := len(rows)
	for c, line := range rows {
		if isEmpty.MatchString(line) {
			continue
		}
		fp.WriteString("      " + line)
		if c < rowCount-1 {
			fp.WriteString(" \\\n")
		} else {
			fp.WriteString("\"\n")
		}
	}
	return nil
}
