proxy-auth-config-downloader
----------------------------

The purpose of this program is to download authentication securely from AWS Secrets Manager and/or SSM.

It then uses this information to generate a configuration file for oauth2proxy.  Future updates will make it
stay resident in the background in order to support secret rotation.

### Why not just use devx-config?

DevXConfig (https://github.com/guardian/devx-config) is a great tool for lifting a generic configuration out of SSM/Secrets
and making it available for scripting. But it felt a bit kludgy to then use bash scripting to turn this into a configuration
for oauth2proxy.

That's not the main reason though; the reason is future features.  The plan for this app is to make it support secret 
rotation, running as a background service on the VM and able to restart oauth2proxy when its secrets change so that underlying
VMs don't need to get regularly restarted.