package main

import (
	"errors"
	"net"

	"github.com/digitalocean/go-metadata"
)

// PrivateInterface returns the network interface name of the provided local
// ip address
func PrivateInterface(ifaces []net.Interface, local string) (string, error) {
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return ``, err
		}

		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPAddr:
				ip := v.IP.String()
				if ip == local {
					return i.Name, nil
				}
			case *net.IPNet:
				ip := v.IP.String()
				if ip == local {
					return i.Name, nil
				}
			}
		}
	}

	return ``, errors.New(`local interface could not be found`)
}

// LocalAddress parses metadata and to find the local private ipv4 interface
// address
func LocalAddress(data *metadata.Metadata) (string, error) {
	privateIface := data.Interfaces[`private`]
	if len(privateIface) >= 1 {
		ipV4 := privateIface[0].IPv4
		if ipV4 == nil {
			return ``, errors.New(`no ipv4 private iface`)
		}

		return ipV4.IPAddress, nil
	}
	return ``, errors.New(`no private interfaces`)
}
