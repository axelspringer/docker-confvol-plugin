[![Taylor Swift](https://img.shields.io/badge/secured%20by-taylor%20swift-brightgreen.svg)](https://twitter.com/SwiftOnSecurity)
[![Volkswagen](https://auchenberg.github.io/volkswagen/volkswargen_ci.svg?v=1)](https://github.com/auchenberg/volkswagen)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# Docker Configuration Volume plugin

## Dependencies

* go
* make
* systemd
* docker

### Install golang dependend packages

```
make deps
```

## Configuration

Main configuration file ```/etc/docker/docker-confvol-plugin```

## Build

Build the whole project
```make build man```
or in short
```make``` 

## Installation

```
sudo -i

make
make install

systemctl daemon-reload
systemctl start docker-confvol-plugin

```

### Use the plugin

You can simply use a direct file mount like this

```--mount volume-driver=confvol,target=/etc/nginx/conf.d/site.conf,source=dev/example/nginx/conf.d/site.conf```

Or the same with a templated configuration

```--mount volume-driver=confvol,target=/etc/nginx/conf.d/site.conf,source=dev/example/nginx/conf.d/site.conf,volume-opt=gen=1```

Or you can mount folders

```--mount volume-driver=confvol,target=/var/www/htdocs/,source=dev/example/nginx/htdocs/```

For the complete example start the vagrant box, the etcd and the etcd browser. Fill the struct from examples/etcd_root to the etcd.

```
docker run \
    --rm \
    --mount volume-driver=confvol,target=/etc/nginx/conf.d/site.conf,source=dev/example/nginx/conf.d/site.conf \ 
    --mount volume-driver=confvol,target=/var/www/htdocs/,source=dev/example/nginx/htdocs/ \
    -p 8080:8080 \
    -d \
    nginx 
```

A more complex example with a templated nginx

```
docker run \
    --rm \
    --mount volume-driver=confvol,target=/etc/nginx/.htpasswd,source=dev/nginx/etc/nginx/.htpasswd,volume-opt=tmpl=1 \    
    --mount volume-driver=confvol,target=/etc/nginx/conf.d/site.conf,source=dev/nginx/etc/nginx/conf.d/site-basicauth.conf \
    --mount volume-driver=confvol,target=/var/www/htdocs/,source=dev/nginx/var/www/htdocs/ \
    -p 8080:8080 \
    -d \
    nginx 
```

## Options

#### Program arguments

* ```--config=<Path>``` Path to the configuration file

#### Docker mount arguments

* ```volume-driver=confvol``` specify the driver  
* ```target=<container-path>``` mount point within the container
* ```source=<conf-path>``` configuration path 
* ```volume-opt=tmpl=1``` evaluated template file
* ```volume-opt=mode=0644``` target file mode bits (in octal)
* ```readonly``` readonly mode 

### Useful ressources

* [https://docs.docker.com/engine/extend/plugin_api/](https://docs.docker.com/engine/extend/plugin_api/)
* [https://docs.docker.com/engine/extend/plugins_volume/
](https://docs.docker.com/engine/extend/plugins_volume/)
* [https://blog.codeship.com/extend-docker-via-plugin/](https://blog.codeship.com/extend-docker-via-plugin/)
* [https://github.com/docker/go-plugins-helpers](https://github.com/docker/go-plugins-helpers)

## License
[Apache-2.0](/LICENSE)
