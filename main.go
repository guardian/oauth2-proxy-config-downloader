package main

import (
	"context"
	"flag"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"log"
	"os"
	"path"
	"strings"
)

/*
WriteString outputs the given string into a textfile in the given directory. Short n simple.
*/
func WriteString(content string, dir string, filename string) error {
	fullPath := path.Join(dir, filename)
	fp, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer fp.Close()

	_, err = fp.WriteString(content)
	return err
}

/*
UpstreamFlag allows us to handle multiple repeated -upstream flags by combining the values into an array
*/
type UpstreamFlag struct {
	Upstreams []string
}

func (u *UpstreamFlag) Set(value string) error {
	u.Upstreams = append(u.Upstreams, value)
	return nil
}

func (u *UpstreamFlag) String() string {
	return strings.Join(u.Upstreams, " ")
}

func main() {
	googleAuthSecret := flag.String("googleAuthJsonSecret", "", "name of the secret in AWS Secrets Manager to obtain the google auth json from")
	googleClientSecret := flag.String("googleAuthClientSecret", "", "name of the secret in AWS Secrets Manager to obtain the google auth client secret")
	sessionCookieSecret := flag.String("sessionCookieSecret", "", "name of the secret in AWS Secrets Manager to obtain the session cookie encoding secret")
	timeoutSeconds := flag.Int("timeout", 5, "maximum number of seconds to wait for a response from the backend service")
	upstreams := &UpstreamFlag{Upstreams: []string{}}
	flag.Var(upstreams, "upstream", "upstream location to forward valid requests on to")
	outputDir := flag.String("out", "/etc/oauth2proxy", "path at which files should be output")
	app := flag.String("app", "", "Application name for generating standardised SSM paths")
	stack := flag.String("stack", "", "Stack name for generating standardised SSM paths")
	stage := flag.String("stage", "", "Stage name for generating standardised SSM paths")
	flag.Parse()

	exitCode := 0

	log.Printf("INFO Upstreams are given as %s", upstreams)
	appId := NewAppId(*app, *stack, *stage)

	if *googleAuthSecret == "" {
		log.Fatal("You must specify a secret name to load in for google json. Use --help for details")
	}
	if *sessionCookieSecret == "" {
		log.Fatal("You must specify a secret name to load in for the session cookie. Use --help for details")
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal("Could not initialise AWS config: ", err)
	}

	secretsClient := secretsmanager.NewFromConfig(cfg)

	googleAuthJson, err := loadGoogleAuthSecretJson(secretsClient, googleAuthSecret, *timeoutSeconds)
	if err != nil {
		log.Printf("ERROR Could not load google auth json from '%s': %s", *googleAuthSecret, err)
		exitCode = 1
	} else {
		if googleAuthJson.Validate() {
			err = googleAuthJson.Output(*outputDir)
			if err != nil {
				log.Printf("ERROR Could not output google auth json to %s: %s", *outputDir, err)
				exitCode = 3
			}
		} else {
			log.Printf("ERROR One or more fields was missing from the google auth JSON. Please ensure that it is copied verbatim.")
			exitCode = 2
		}
	}

	googleClientValue, err := loadSecretString(secretsClient, googleClientSecret, *timeoutSeconds)
	if err != nil {
		log.Printf("ERROR Could not load google auth client secret from '%s': %s", *googleClientSecret, err)
		exitCode = 1
	}

	sessionCookie, err := loadSecretString(secretsClient, sessionCookieSecret, *timeoutSeconds)
	if err != nil {
		log.Printf("ERROR Could not load session cookie from '%s': %s", *sessionCookieSecret, err)
		exitCode = 1
	}

	ssmClient := ssm.NewFromConfig(cfg)
	adminEmail, err := GetSSMValueString(ssmClient, appId, "googleAdminEmail", *timeoutSeconds)
	if err != nil {
		log.Printf("ERROR Could not load google admin email: %s", err)
		exitCode = 1
	}
	allowedGroupList, err := GetSSMValueStringList(ssmClient, appId, "googleAuthAllowedGroups", *timeoutSeconds)
	if err != nil {
		log.Printf("ERROR Could not load google groups as string list: %s", err)
		exitCode = 1
	}
	emailDomain, err := GetSSMValueString(ssmClient, appId, "googleAuthAllowedEmailDomain", *timeoutSeconds)
	if err != nil {
		log.Printf("ERROR Could not load google allowed email address domain: %s", err)
		exitCode = 1
	}
	appDomain, err := GetSSMValueString(ssmClient, appId, "appDomainName", *timeoutSeconds)
	if err != nil {
		log.Printf("ERROR Could not load domain name of the protected app: %s", err)
		exitCode = 1
	}
	googleClientId, err := GetSSMValueString(ssmClient, appId, "googleClientId", *timeoutSeconds)

	if exitCode != 0 {
		os.Exit(exitCode)
	}

	genConfig := NewOAuth2ProxyConfig(googleAuthJson, googleClientId, googleClientValue, adminEmail, allowedGroupList, emailDomain, sessionCookie, appId.CookieName(), appDomain, upstreams.Upstreams, *outputDir)

	err = genConfig.Generate(*outputDir)
	if err != nil {
		log.Print("ERROR Could not generate config file: ", err)
		os.Exit(3)
	}
	os.Exit(exitCode)
}
