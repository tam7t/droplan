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

	accessToken := os.Getenv("DO_KEY")
	if accessToken == "" {
		log.Fatal("Usage: DO_KEY environment variable must be set.")
	}

	peerTag := os.Getenv("DO_TAG")

	// PUBLIC=true will tell us to block traffic on the public interface
	public := os.Getenv("PUBLIC")

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

	// collect local network interface information
	ifaces, err := net.Interfaces()
	failIfErr(err)

	pubAddr, err := PublicAddress(mData)
	failIfErr(err)

	if public == "true" {
		publicPeers := PublicDroplets(drops)

		// find public iface name
		iface, err := FindInterfaceName(ifaces, pubAddr)
		failIfErr(err)

		// setup droplan-peers-public chain for public interface
		err = Setup(ipt, iface, "droplan-peers-public")
		failIfErr(err)

		// update droplan-peers-public
		err = UpdatePeers(ipt, publicPeers, "droplan-peers-public")
		failIfErr(err)
		log.Printf("Added %d peers to droplan-peers-public", len(publicPeers))
	}

	privAddr, err := PrivateAddress(mData)
	failIfErr(err)

	privatePeers, ok := SortDroplets(drops)[region]
	if !ok {
		log.Printf("No droplets listed in region [%s]", region)
	}

	// find private iface name
	iface, err := FindInterfaceName(ifaces, privAddr)
	if public != "" && err != nil && err.Error() == "no private interfaces" {
		os.Exit(0)
	}
	failIfErr(err)

	// setup droplan-peers chain for private interface
	err = Setup(ipt, iface, "droplan-peers")
	failIfErr(err)

	// update droplan-peers
	err = UpdatePeers(ipt, privatePeers, "droplan-peers")
	failIfErr(err)
	log.Printf("Added %d peers to droplan-peers", len(privatePeers))
}

func failIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
