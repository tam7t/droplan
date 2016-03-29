# droplan [![Build Status](http://img.shields.io/travis/tam7t/droplan.svg?style=flat-square)](https://travis-ci.org/tam7t/droplan) [![Gitter](https://img.shields.io/gitter/room/tam7t/droplan.js.svg?style=flat-square)](https://gitter.im/tam7t/droplan)

## About

This utility helps secure the `private` interface on DigitalOcean droplets by
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

## Development

### Dependencies

Dependencies are vendored with [govendor](https://github.com/kardianos/govendor).

### Build

A `Makefile` is included:
  * `test` - runs unit tests
  * `build` - builds `droplan` on the current platform
  * `release` - builds releasable artifacts
