package main

import (
	"flag"
	"log"
	"net"
	"os"

	"github.com/coreos/go-iptables/iptables"
	"github.com/digitalocean/go-metadata"
	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

var appVersion string

func main() {
	version := flag.Bool("version", false, "Print the version and exit.")
	flag.Parse()
	if *version {
		log.Printf(appVersion)
		os.Exit(0)
	}

	accessToken := os.Getenv(`DO_KEY`)
	if accessToken == `` {
		log.Fatal(`Usage: DO_KEY environment variable must be set.`)
	}

	peerTag := os.Getenv(`DO_TAG`)

	// setup dependencies
	oauthClient := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken}))
	apiClient := godo.NewClient(oauthClient)
	metaClient := metadata.NewClient()
	ipt, err := iptables.New()
	failIfErr(err)

	// collect needed metadata from metadata service
	region, err := metaClient.Region()
	failIfErr(err)
	mData, err := metaClient.Metadata()
	failIfErr(err)

	// collect list of all droplets
	var drops []godo.Droplet
	if peerTag != "" {
		drops, err = DropletListTags(apiClient.Droplets, peerTag)
	} else {
		drops, err = DropletList(apiClient.Droplets)
	}
	failIfErr(err)

	allowed, ok := SortDroplets(drops)[region]
	if !ok {
		log.Fatalf(`No droplets listed in region [%s]`, region)
	}

	// collect local network interface information
	local, err := LocalAddress(mData)
	failIfErr(err)
	ifaces, err := net.Interfaces()
	failIfErr(err)
	iface, err := PrivateInterface(ifaces, local)
	failIfErr(err)

	// setup droplan-peers chain for local interface
	err = Setup(ipt, iface)
	failIfErr(err)

	// update droplan-peers
	err = UpdatePeers(ipt, allowed)
	failIfErr(err)
	log.Printf(`Added %d peers to droplan-peers`, len(allowed))
}

func failIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
