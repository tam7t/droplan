package main

import (
	"errors"
	"net/url"
	"testing"

	"github.com/digitalocean/godo"
	. "github.com/franela/goblin"
)

func TestPeers(t *testing.T) {
	g := Goblin(t)

	g.Describe(`SortDroplets`, func() {
		region := &godo.Region{Slug: `nyc1`}
		droplet := godo.Droplet{Region: region}
		var output map[string][]string

		g.BeforeEach(func() {
			output = SortDroplets([]godo.Droplet{droplet})
		})

		g.Describe(`without a private network`, func() {
			g.Before(func() {
				droplet.Networks = &godo.Networks{
					V4: []godo.NetworkV4{
						godo.NetworkV4{IPAddress: `192.168.0.0`, Type: `public`},
					},
				}
			})

			g.It(`is not included in the output`, func() {
				_, exists := output[region.Slug]
				g.Assert(exists).Equal(false)
			})
		})

		g.Describe(`with a private network`, func() {
			g.Before(func() {
				droplet.Networks = &godo.Networks{
					V4: []godo.NetworkV4{
						godo.NetworkV4{IPAddress: `192.168.0.0`, Type: `private`},
					},
				}
			})

			g.It(`is included in the output`, func() {
				g.Assert(output[region.Slug]).Equal([]string{`192.168.0.0`})
			})
		})
	})

	g.Describe(`DropletList`, func() {
		var sds *stubDropletService
		g.BeforeEach(func() {
			sds = &stubDropletService{}
		})

		g.Describe(`with no droplets`, func() {
			g.BeforeEach(func() {
				sds.list = func(a *godo.ListOptions) ([]godo.Droplet, *godo.Response, error) {
					resp := &godo.Response{}
					resp.Links = nil
					return []godo.Droplet{}, resp, nil
				}
			})

			g.It(`returns no droplets`, func() {
				drops, err := DropletList(sds)
				g.Assert(drops).Equal([]godo.Droplet{})
				g.Assert(err).Equal(nil)
			})
		})

		g.Describe(`with a single page of droplets`, func() {
			g.BeforeEach(func() {
				sds.list = func(a *godo.ListOptions) ([]godo.Droplet, *godo.Response, error) {
					resp := &godo.Response{}
					resp.Links = nil
					return []godo.Droplet{{Name: `foobar`}}, resp, nil
				}
			})

			g.It(`returns the droplets`, func() {
				drops, err := DropletList(sds)
				g.Assert(drops).Equal([]godo.Droplet{{Name: `foobar`}})
				g.Assert(err).Equal(nil)
			})
		})

		g.Describe(`with a multiple pages of droplets`, func() {
			g.BeforeEach(func() {
				sds.list = func(a *godo.ListOptions) ([]godo.Droplet, *godo.Response, error) {
					resp := &godo.Response{}
					drops := []godo.Droplet{}
					if a.Page == 0 {
						resp.Links = &godo.Links{Pages: &godo.Pages{Next: `http://example.com/droplets?page=2`, Last: `http://example.com/droplets?page=2`}}
						drops = append(drops, godo.Droplet{Name: `firstPage`})
					} else {
						resp.Links = &godo.Links{Pages: &godo.Pages{Prev: `http://example.com/droplets?page=1`}}
						drops = append(drops, godo.Droplet{Name: `secondPage`})
					}
					return drops, resp, nil
				}
			})

			g.It(`returns the droplets`, func() {
				drops, err := DropletList(sds)
				g.Assert(drops).Equal([]godo.Droplet{{Name: `firstPage`}, {Name: `secondPage`}})
				g.Assert(err).Equal(nil)
			})
		})

		g.Describe("when droplet services list errors", func() {
			g.BeforeEach(func() {
				sds.list = func(a *godo.ListOptions) ([]godo.Droplet, *godo.Response, error) {
					return []godo.Droplet{}, nil, errors.New("asdf")
				}
			})

			g.It("returns an error", func() {
				_, err := DropletList(sds)
				g.Assert(err).Equal(errors.New("asdf"))
			})
		})

		g.Describe("when current page errors", func() {
			g.BeforeEach(func() {
				sds.list = func(a *godo.ListOptions) ([]godo.Droplet, *godo.Response, error) {
					resp := &godo.Response{}
					resp.Links = &godo.Links{Pages: &godo.Pages{Prev: "page=)", Last: "page="}}
					return []godo.Droplet{{Name: "foobar"}}, resp, nil
				}
			})

			g.It("returns an error", func() {
				_, err := DropletList(sds)

				g.Assert(err.(*url.Error).Op).Equal("parse")
			})
		})
	})
}

type stubDropletService struct {
	list           func(*godo.ListOptions) ([]godo.Droplet, *godo.Response, error)
	listTag        func(string, *godo.ListOptions) ([]godo.Droplet, *godo.Response, error)
	get            func(int) (*godo.Droplet, *godo.Response, error)
	create         func(*godo.DropletCreateRequest) (*godo.Droplet, *godo.Response, error)
	createMultiple func(*godo.DropletMultiCreateRequest) ([]godo.Droplet, *godo.Response, error)
	delete         func(int) (*godo.Response, error)
	deleteTag      func(string) (*godo.Response, error)
	kernels        func(int, *godo.ListOptions) ([]godo.Kernel, *godo.Response, error)
	snapshots      func(int, *godo.ListOptions) ([]godo.Image, *godo.Response, error)
	backups        func(int, *godo.ListOptions) ([]godo.Image, *godo.Response, error)
	actions        func(int, *godo.ListOptions) ([]godo.Action, *godo.Response, error)
	neighbors      func(int) ([]godo.Droplet, *godo.Response, error)
}

func (sds *stubDropletService) List(a *godo.ListOptions) ([]godo.Droplet, *godo.Response, error) {
	return sds.list(a)
}

func (sds *stubDropletService) ListByTag(a string, b *godo.ListOptions) ([]godo.Droplet, *godo.Response, error) {
	return sds.listTag(a, b)
}

func (sds *stubDropletService) Get(a int) (*godo.Droplet, *godo.Response, error) {
	return sds.get(a)
}

func (sds *stubDropletService) Create(a *godo.DropletCreateRequest) (*godo.Droplet, *godo.Response, error) {
	return sds.create(a)
}

func (sds *stubDropletService) CreateMultiple(a *godo.DropletMultiCreateRequest) ([]godo.Droplet, *godo.Response, error) {
	return sds.createMultiple(a)
}

func (sds *stubDropletService) Delete(a int) (*godo.Response, error) {
	return sds.delete(a)
}

func (sds *stubDropletService) DeleteByTag(a string) (*godo.Response, error) {
	return sds.deleteTag(a)
}

func (sds *stubDropletService) Kernels(a int, b *godo.ListOptions) ([]godo.Kernel, *godo.Response, error) {
	return sds.kernels(a, b)
}

func (sds *stubDropletService) Snapshots(a int, b *godo.ListOptions) ([]godo.Image, *godo.Response, error) {
	return sds.snapshots(a, b)
}

func (sds *stubDropletService) Backups(a int, b *godo.ListOptions) ([]godo.Image, *godo.Response, error) {
	return sds.backups(a, b)
}

func (sds *stubDropletService) Actions(a int, b *godo.ListOptions) ([]godo.Action, *godo.Response, error) {
	return sds.actions(a, b)
}

func (sds *stubDropletService) Neighbors(a int) ([]godo.Droplet, *godo.Response, error) {
	return sds.neighbors(a)
}
