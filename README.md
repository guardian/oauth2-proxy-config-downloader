proxy-auth-config-downloader
----------------------------

The purpose of this program is to download authentication securely from AWS Secrets Manager and/or SSM.

It then uses this information to generate a configuration file for oauth2proxy.  Future updates will make it
stay resident in the background in order to support secret rotation.

