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

```
docker run --volume-driver confvol --volume confvol-48E014D6-1B7F-4634-B883-3B787AC84032:/data alpine sleep 60
```

```
docker run \
    --mount volume-driver=confvol,source=48E014D6-1B7F-4634-B883-3B787AC84032,target=/var/www/htdocs,volume-opt=o=/app/nginx/htdocs \
    --mount volume-driver=confvol,source=48E014D6-1B7F-4634-B883-3B787AC84032,target=/etc/nginx/conf.d/default,volume-opt=o=/app/nginx/conf \
    alpine sleep 60
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
