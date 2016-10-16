package main

import "github.com/digitalocean/godo"

// SortDroplets returns a map (keyed by region slug) of droplets with private ip
// interfaces
func SortDroplets(droplets []godo.Droplet) map[string][]string {
	netDrops := map[string][]string{}

	for _, droplet := range droplets {
		for _, net := range droplet.Networks.V4 {
			if net.Type == "private" {
				_, ok := netDrops[droplet.Region.Slug]
				if ok {
					netDrops[droplet.Region.Slug] = append(netDrops[droplet.Region.Slug], net.IPAddress)
				} else {
					netDrops[droplet.Region.Slug] = []string{net.IPAddress}
				}
			}
		}
	}

	return netDrops
}

// PublicDroplets returns an array of all the public ip interfaces of the provided
// droplets
func PublicDroplets(droplets []godo.Droplet) []string {
	netDrops := []string{}
	for _, droplet := range droplets {
		for _, net := range droplet.Networks.V4 {
			if net.Type == "public" {
				netDrops = append(netDrops, net.IPAddress)
			}
		}
	}
	return netDrops
}

// DropletList paginates through the digitalocean API to return a list of all
// droplets
func DropletList(ds godo.DropletsService) ([]godo.Droplet, error) {
	// create a list to hold our droplets
	list := []godo.Droplet{}

	// create options. initially, these will be blank
	opt := &godo.ListOptions{}
	for {
		droplets, resp, err := ds.List(opt)

		if err != nil {
			return nil, err
		}

		// append the current page's droplets to our list
		for _, d := range droplets {
			list = append(list, d)
		}

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}

	return list, nil
}

// DropletListTags paginates through the digitalocean API to return a list of
// all droplets with the given tag
func DropletListTags(ds godo.DropletsService, tag string) ([]godo.Droplet, error) {
	// create a list to hold our droplets
	list := []godo.Droplet{}

	// create options. initially, these will be blank
	opt := &godo.ListOptions{}
	for {
		droplets, resp, err := ds.ListByTag(tag, opt)

		if err != nil {
			return nil, err
		}

		// append the current page's droplets to our list
		for _, d := range droplets {
			list = append(list, d)
		}

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}

	return list, nil
}
