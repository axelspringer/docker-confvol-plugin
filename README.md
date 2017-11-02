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

You can use a direct file mount like

```--mount volume-driver=confvol,target=/etc/nginx/conf.d/site.conf,source=/dev/example/nginx/conf.d/site.conf```

Or you can mount folders

```--mount volume-driver=confvol,target=/var/www/htdocs/,source=/dev/example/nginx/htdocs/```

For the complete example start the vagrant box, the etcd and the etcd browser. Fill the struct from examples/etcd_root to the etcd.
Then run ... 

```
docker run \
    --rm \
    --mount volume-driver=confvol,target=/etc/nginx/conf.d/site.conf,source=/dev/example/nginx/conf.d/site.conf \ 
    --mount volume-driver=confvol,target=/var/www/htdocs/,source=/dev/example/nginx/htdocs/ \
    -p 8080:8080 \
    -d \
    nginx 
```

## Docs

### Useful ressources

* [https://docs.docker.com/engine/extend/plugin_api/](https://docs.docker.com/engine/extend/plugin_api/)
* [https://docs.docker.com/engine/extend/plugins_volume/
](https://docs.docker.com/engine/extend/plugins_volume/)
* [https://blog.codeship.com/extend-docker-via-plugin/](https://blog.codeship.com/extend-docker-via-plugin/)
* [https://github.com/docker/go-plugins-helpers](https://github.com/docker/go-plugins-helpers)

## License
[Apache-2.0](/LICENSE)
