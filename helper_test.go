package main

import (
	"encoding/json"
	"errors"
	"net"
	"reflect"
	"strings"
	"testing"

	"github.com/digitalocean/go-metadata"
)

func TestFindInterfaceName(t *testing.T) {
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

	tests := []struct {
		name   string
		ifaces []net.Interface
		local  string
		exp    string
		expErr error
	}{
		{
			name:   "valid ip",
			ifaces: ifaces,
			local:  ip,
			exp:    ifaces[0].Name,
			expErr: nil,
		},
		{
			name:   "invalid ip",
			ifaces: ifaces,
			local:  "somethingbad",
			exp:    "",
			expErr: errors.New("local interface could not be found"),
		},
	}

	for _, test := range tests {
		out, err := FindInterfaceName(test.ifaces, test.local)
		if !reflect.DeepEqual(err, test.expErr) {
			t.Logf("want:%v", test.expErr)
			t.Logf("got:%v", err)
			t.Fatalf("test case failed: %s", test.name)
		}
		if out != test.exp {
			t.Logf("want:%v", test.exp)
			t.Logf("got:%v", out)
			t.Fatalf("test case failed: %s", test.name)
		}
	}
}

func TestPrivateAddress(t *testing.T) {
	tests := []struct {
		name   string
		data   *metadata.Metadata
		exp    string
		expErr error
	}{
		{
			name:   "private ipv4 address",
			data:   decodeMetadata(`{"interfaces": {"private": [{"ipv4": {"ip_address": "privateIP"}}]}}`),
			exp:    "privateIP",
			expErr: nil,
		},
		{
			name:   "private ipv6 address",
			data:   decodeMetadata(`{"interfaces": {"private": [{"ipv6": {"ip_address": "privateIP"}}]}}`),
			exp:    "",
			expErr: errors.New("no ipv4 private iface"),
		},
		{
			name:   "no private addresses",
			data:   &metadata.Metadata{},
			exp:    "",
			expErr: errors.New("no private interfaces"),
		},
	}

	for _, test := range tests {
		out, err := PrivateAddress(test.data)
		if !reflect.DeepEqual(err, test.expErr) {
			t.Logf("want:%v", test.expErr)
			t.Logf("got:%v", err)
			t.Fatalf("test case failed: %s", test.name)
		}
		if out != test.exp {
			t.Logf("want:%v", test.exp)
			t.Logf("got:%v", out)
			t.Fatalf("test case failed: %s", test.name)
		}
	}
}

func TestPublicAddress(t *testing.T) {
	tests := []struct {
		name   string
		data   *metadata.Metadata
		exp    string
		expErr error
	}{
		{
			name:   "public ipv4 address",
			data:   decodeMetadata(`{"interfaces": {"public": [{"ipv4": {"ip_address": "publicIP"}}]}}`),
			exp:    "publicIP",
			expErr: nil,
		},
		{
			name:   "public ipv6 address",
			data:   decodeMetadata(`{"interfaces": {"public": [{"ipv6": {"ip_address": "publicIP"}}]}}`),
			exp:    "",
			expErr: errors.New("no ipv4 public iface"),
		},
		{
			name:   "no public addresses",
			data:   &metadata.Metadata{},
			exp:    "",
			expErr: errors.New("no public interfaces"),
		},
	}

	for _, test := range tests {
		out, err := PublicAddress(test.data)
		if !reflect.DeepEqual(err, test.expErr) {
			t.Logf("want:%v", test.expErr)
			t.Logf("got:%v", err)
			t.Fatalf("test case failed: %s", test.name)
		}
		if out != test.exp {
			t.Logf("want:%v", test.exp)
			t.Logf("got:%v", out)
			t.Fatalf("test case failed: %s", test.name)
		}
	}
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
