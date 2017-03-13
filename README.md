# droplan [![Build Status](http://img.shields.io/travis/tam7t/droplan.svg?style=flat-square)](https://travis-ci.org/tam7t/droplan) [![Gitter](https://img.shields.io/gitter/room/tam7t/droplan.js.svg?style=flat-square)](https://gitter.im/tam7t/droplan)

## About

This utility helps secure the network interfaces on DigitalOcean droplets by
adding `iptable` rules that only allow traffic from your other droplets. `droplan`
queries the DigitalOcean API and automatically updates `iptable` rules.

## Installation

The latest release is available on the github [release page](https://github.com/tam7t/droplan/releases).

You can setup a cron job to run every 5 minutes in `/etc/cron.d`

```
*/5 * * * * root PATH=/sbin DO_KEY=READONLY_KEY /usr/local/bin/droplan >/var/log/droplan.log 2>&1
```

## Usage

```
DO_KEY=<read_only_api_token> /path/to/droplan
```

The `iptables` rules added by `droplan` are equivalent to:

```
-N droplan-peers # create a new chain
-A INPUT -i eth1 -j droplan-peers # add chain to private interface
-A INPUT -i eth1 -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
-A INPUT -i eth1 -j DROP # add default DROP rule to private interface
-A droplan-peers -s <PEER>/32 -j ACCEPT # allow traffic from PEER ip address
```

### Tags
Access can be limited to a subset of droplets using [tags](https://developers.digitalocean.com/documentation/v2/#tags).
The `DO_TAG` environment variable tells `droplan` to only allow access to
droplets with the specified tag.

### Public Interface
Add the `PUBLIC=true` environment variable and `droplan` will maintain an
iptables chain of `droplan-peers-public` with the public ip addresses of
peers and add a default drop rule to the `eth0` interface.

**NOTE:** This will prevent you from being able to directly ssh into your droplet.

## Development

### Dependencies

Dependencies are vendored with [govendor](https://github.com/kardianos/govendor).

### Build

A `Makefile` is included:
  * `test` - runs unit tests
  * `build` - builds `droplan` on the current platform
  * `release` - builds releasable artifacts


## Docker image:

There's a prebuilt [docker image][1] (currently it's third party, but I hope it'll get merged back to tam7t/droplan).

Example usage:

```sh
docker run -d --restart=always --net=host -e DO_KEY=$your_digitalocean_api_key tam7t/droplan
```

- `-d --restart=always` starts the container in the background and restarts it on error
- `--net=host` is required because we want to affect the host's firewall rules, not the container's
- you have to specify your DigitalOcean API key (using `-e DO_KEY`)
- you can add `-e PUBLIC=true` or `-e DO_TAG=tagname` as described above
- specify `-e DO_INTERVAL=120` to change the delay between droplan invocations (default: '300' (5 minutes) ) 


[1]: https://hub.docker.com/r/tam7t/droplan/
