% DOCKER-CONFVOL-PLUGIN(8) 
% Jan Michalowsky 
% NOVEMBER 2017
# NAME
docker-confvol-plugin - Docker conf volume driver for libkv backends

# SYNOPSIS
**docker-confvol-plugin**
[**-debug**]
[**-version**]

# STATE
Experimental! Do not use in production!!!

# DESCRIPTION
This plugin can be used to create volumes that based on a kv tree

# USAGE
Start the docker daemon before starting the docker-confvol-plugin daemon. 
You can start docker daemon using command:
```bash
systemctl start docker 
```
Once docker daemon is up and running, you can start docker-confvol-plugin daemon
using command:
```bash
systemctl start docker-confvol-plugin
``` 

# OPTIONS

# EXAMPLES
