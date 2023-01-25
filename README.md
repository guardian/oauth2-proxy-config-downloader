proxy-auth-config-downloader
----------------------------

The purpose of this program is to download authentication securely from AWS Secrets Manager and/or SSM.

It then uses this information to generate a configuration file for oauth2proxy.  Future updates will make it
stay resident in the background in order to support secret rotation.

## How do I obtain it?

Either use `go install`, clone this repo and build it with `go build` or take the precompiled binary from the
"releases" page.

If you want to get hold of the latest version to build it into another project's deliverable, you can simply use
a command like 

```bash
curl -L https://github.com/guardian/oauth2-proxy-config-downloader/releases/latest/download/oauth2-proxy-config-downloader.aarch64.gz | gunzip > oauth2-proxy-config-downloader 
```

to download the latest ARM architecture build.  Substitute the filename for any other filename from the release page to get
that file (note that it uses redirects internally, so you need `-L` to tell curl to follow them)


## How does it work?

You simply run the command with the appropriate arguments:
```
./proxy-auth-config-downloader -help
Usage of ./proxy-auth-config-downloader:
  -app string
    	Application name for generating standardised SSM paths
  -googleAuthClientSecret string
    	name of the secret in AWS Secrets Manager to obtain the google auth client secret
  -googleAuthJsonSecret string
    	name of the secret in AWS Secrets Manager to obtain the google auth json from
  -out string
    	path at which files should be output (default "/etc/oauth2proxy")
  -sessionCookieSecret string
    	name of the secret in AWS Secrets Manager to obtain the session cookie encoding secret
  -stack string
    	Stack name for generating standardised SSM paths
  -stage string
    	Stage name for generating standardised SSM paths
  -timeout int
    	maximum number of seconds to wait for a response from the backend service (default 5)
  -upstream value
    	upstream location to forward valid requests on to
```

For example:
```
AWS_REGION=eu-west-1 proxy-auth-config-downloader -googleAuthJsonSecret /${Stage}/${Stack}/${App}/GoogleServiceAcct \
        -sessionCookieSecret /${Stage}/${Stack}/${App}/CookieGenSecret \
        -app ${App} \
        -stack ${Stack} \
        -stage ${Stage} \
        -googleAuthClientSecret /${Stage}/${Stack}/${App}/GoogleAuthClientSecret \
        -upstream https://${UpstreamServiceHost}/${UpstreamServiceSubpath} \
        -upstream file:///var/www/#/
```

(where `${}` denotes an environment variable or other substitution)

This will output the following files into `/etc/oauth2proxy`:
- `oauth2proxy.env` - an environment file, suitable for loading into systemd via an `EnvironmentFile=` directive.  This gives
    you a single variable called `$OPTS` which is the commandline options for oauth2proxy
- `auth.json` - json format file giving the google credentials for authenticaion. This is required by oauth2proxy to make
    google auth work.

### Requirements

In order to work, it expects the following keys to be present in SSM and for the user (or EC2 role) to have permission
to view / decrypt them:
- `/${Stage}/${Stack}/${App}/googleAdminEmail` (type String) - name of a Google Workspace administrator account that the service account
is allowed to "impersonate" in order to do group-membership checks
- `/${Stage}/${Stack}/${App}/googleAuthAllowedGroups` (type StringList) - list of google workspace groups. A user must belong to
at least one of these or they will not be allowed access.
- `/${Stage}/${Stack}/${App}/googleAuthAllowedEmailDomain` (type String) - Google Workspace email domain that a user must belong to
before they will be granted access
- `/${Stage}/${Stack}/${App}/appDomainName` (type String) - DNS domain name of the app that is being protected (used for generating
oauth2 callback URLs)
- `/${Stage}/${Stack}/${App}/googleClientId` (type String) - ID of the google oauth2 client

Additionally, it expects the following secrets to be present in Secrets Manager and for the user (or EC2 role) to have
permission to view / decrypt them:

- `/${Stage}/${Stack}/${App}/GoogleServiceAcct` (json file) - full JSON of a Google Service Account. This is used for
performing authentication and group checks (note, the full name of the secret is customisable from the commandline)
- `/${Stage}/${Stack}/${App}/CookieGenSecret` (string) - randomised string used as a salt for generating secure cookies
- `/${Stage}/${Stack}/${App}/GoogleAuthClientSecret` (string) - secret key provided with the google authentication credentials
in order to allow authentication to proceed

### Loading

Given you have the oauth2-proxy binary installed in `/usr/local/bin`,  the following .service file will allow you to load
it in with the given configuration:

```unit file (systemd)
[Unit]
Description=OAuth2 Proxy
Documentation=https://github.com/oauth2-proxy/oauth2-proxy
After=syslog.target network.target remote-fs.target nss-lookup.target

[Service]
Type=simple
EnvironmentFile=/etc/oauth2proxy/oauth2proxy.env
ExecStart=/usr/local/bin/oauth2-proxy $OPTS
User=oauth2proxy

[Install]
WantedBy=multi-user.target
```

Note that this will error if the config downloader has not been run, because `/etc/oauth2proxy/oauth2proxy.env` will not
have been created yet.

### Why not just use devx-config?

DevXConfig (https://github.com/guardian/devx-config) is a great tool for lifting a generic configuration out of SSM/Secrets
and making it available for scripting. But it felt a bit kludgy to then use bash scripting to turn this into a configuration
for oauth2proxy.

That's not the main reason though; the reason is future features.  The plan for this app is to make it support secret 
rotation, running as a background service on the VM and able to reload oauth2proxy when its secrets change so that underlying
VMs don't need to get regularly restarted.