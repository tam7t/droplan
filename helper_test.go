package main

import (
	"encoding/json"
	"errors"
	"net"
	"strings"
	"testing"

	"github.com/digitalocean/go-metadata"
	. "github.com/franela/goblin"
)

func TestHelper(t *testing.T) {
	g := Goblin(t)

	g.Describe(`PrivateInterface`, func() {
		g.Describe(`with a valid interface`, func() {
			g.It(`returns the correct interface`, func() {
				ifaces, _ := net.Interfaces()
				addrs, _ := ifaces[0].Addrs()
				localAddr := addrs[0]
				ip := localAddr.String()
				switch v := localAddr.(type) {
				case *net.IPAddr:
					ip = v.IP.String()
				case *net.IPNet:
					ip = v.IP.String()
				}
				ifaceName, err := PrivateInterface(ifaces, ip)
				g.Assert(ifaceName).Equal(ifaces[0].Name)
				g.Assert(err).Equal(nil)
			})
		})

		g.It(`with an invalid IP`, func() {
			ifaces, _ := net.Interfaces()
			ifaceName, err := PrivateInterface(ifaces, `somethingBad`)
			g.Assert(ifaceName).Equal(``)
			g.Assert(err).Equal(errors.New(`local interface could not be found`))
		})
	})

	g.Describe(`LocalAddress`, func() {
		g.Describe(`with a private interface`, func() {
			g.Describe(`with an ipv4 address`, func() {
				g.It(`returns the ip address`, func() {
					data := decodeMetadata(`{"interfaces": {"private": [{"ipv4": {"ip_address": "privateIP"}}]}}`)
					addr, _ := LocalAddress(data)
					g.Assert(addr).Equal(`privateIP`)
				})
			})

			g.Describe(`without an ipv4 address`, func() {
				g.It(`returns an error`, func() {
					data := decodeMetadata(`{"interfaces": {"public": [{"ipv4": {"ip_address": "publicIP"}}]}}`)
					_, err := LocalAddress(data)
					g.Assert(err).Equal(errors.New(`no private interfaces`))
				})
			})
		})

		g.Describe(`without a private interface`, func() {
			g.It(`returns an error`, func() {
				data := &metadata.Metadata{}
				_, err := LocalAddress(data)
				g.Assert(err).Equal(errors.New(`no private interfaces`))
			})
		})
	})
}

func decodeMetadata(data string) *metadata.Metadata {
	var output metadata.Metadata
	var err error

	err = json.NewDecoder(strings.NewReader(data)).Decode(&output)
	if err != nil {
		panic(err)
	}
	return &output
}
