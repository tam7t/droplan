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
