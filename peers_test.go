package main

import (
	"errors"
	"reflect"
	"testing"

	"github.com/digitalocean/godo"
)

func TestDropletList(t *testing.T) {
	tests := []struct {
		name             string
		ds               *stubDropletService
		expectedDroplets []godo.Droplet
		expectedError    error
	}{
		{
			name: "no droplets",
			ds: &stubDropletService{
				list: func(b *godo.ListOptions) ([]godo.Droplet, *godo.Response, error) {
					resp := &godo.Response{}
					resp.Links = nil
					return []godo.Droplet{}, resp, nil
				},
			},
			expectedDroplets: []godo.Droplet{},
		},
		{
			name: "single page of droplets",
			ds: &stubDropletService{
				list: func(b *godo.ListOptions) ([]godo.Droplet, *godo.Response, error) {
					resp := &godo.Response{}
					resp.Links = nil
					return []godo.Droplet{{Name: "foobar"}}, resp, nil
				},
			},
			expectedDroplets: []godo.Droplet{{Name: "foobar"}},
		},
		{
			name: "multiple pages of droplets",
			ds: &stubDropletService{
				list: func(b *godo.ListOptions) ([]godo.Droplet, *godo.Response, error) {
					resp := &godo.Response{}
					drops := []godo.Droplet{}
					if b.Page == 0 {
						resp.Links = &godo.Links{Pages: &godo.Pages{Next: "http://example.com/droplets?page=2", Last: "http://example.com/droplets?page=2"}}
						drops = append(drops, godo.Droplet{Name: "firstPage"})
					} else {
						resp.Links = &godo.Links{Pages: &godo.Pages{Prev: "http://example.com/droplets?page=1"}}
						drops = append(drops, godo.Droplet{Name: "secondPage"})
					}
					return drops, resp, nil
				},
			},
			expectedDroplets: []godo.Droplet{{Name: "firstPage"}, {Name: "secondPage"}},
		},
		{
			name: "list errors",
			ds: &stubDropletService{
				list: func(b *godo.ListOptions) ([]godo.Droplet, *godo.Response, error) {
					return []godo.Droplet{}, nil, errors.New("asdf")
				},
			},
			expectedError: errors.New("asdf"),
		},
		{
			name: "current page errors",
			ds: &stubDropletService{
				list: func(b *godo.ListOptions) ([]godo.Droplet, *godo.Response, error) {
					resp := &godo.Response{}
					resp.Links = &godo.Links{Pages: &godo.Pages{Prev: "page=)", Last: "page="}}
					return []godo.Droplet{{Name: "foobar"}}, resp, nil
				},
			},
			expectedError: errors.New("parse page=): invalid URI for request"),
		},
	}

	for _, test := range tests {
		out, err := DropletList(test.ds)
		if !reflect.DeepEqual(err, test.expectedError) {
			if err.Error() != test.expectedError.Error() {
				t.Logf("want:%v", test.expectedError)
				t.Logf("got:%v", err)
				t.Fatalf("test case failed: %s", test.name)
			}
		}
		if !reflect.DeepEqual(out, test.expectedDroplets) {
			t.Logf("want:%v", test.expectedDroplets)
			t.Logf("got:%v", out)
			t.Fatalf("test case failed: %s", test.name)
		}
	}
}

func TestDropletListTags(t *testing.T) {
	tests := []struct {
		name             string
		ds               *stubDropletService
		expectedDroplets []godo.Droplet
		expectedError    error
	}{
		{
			name: "no droplets",
			ds: &stubDropletService{
				listTag: func(a string, b *godo.ListOptions) ([]godo.Droplet, *godo.Response, error) {
					resp := &godo.Response{}
					resp.Links = nil
					return []godo.Droplet{}, resp, nil
				},
			},
			expectedDroplets: []godo.Droplet{},
		},
		{
			name: "single page of droplets",
			ds: &stubDropletService{
				listTag: func(a string, b *godo.ListOptions) ([]godo.Droplet, *godo.Response, error) {
					resp := &godo.Response{}
					resp.Links = nil
					return []godo.Droplet{{Name: "foobar"}}, resp, nil
				},
			},
			expectedDroplets: []godo.Droplet{{Name: "foobar"}},
		},
		{
			name: "multiple pages of droplets",
			ds: &stubDropletService{
				listTag: func(a string, b *godo.ListOptions) ([]godo.Droplet, *godo.Response, error) {
					resp := &godo.Response{}
					drops := []godo.Droplet{}
					if b.Page == 0 {
						resp.Links = &godo.Links{Pages: &godo.Pages{Next: "http://example.com/droplets?page=2", Last: "http://example.com/droplets?page=2"}}
						drops = append(drops, godo.Droplet{Name: "firstPage"})
					} else {
						resp.Links = &godo.Links{Pages: &godo.Pages{Prev: "http://example.com/droplets?page=1"}}
						drops = append(drops, godo.Droplet{Name: "secondPage"})
					}
					return drops, resp, nil
				},
			},
			expectedDroplets: []godo.Droplet{{Name: "firstPage"}, {Name: "secondPage"}},
		},
		{
			name: "list errors",
			ds: &stubDropletService{
				listTag: func(a string, b *godo.ListOptions) ([]godo.Droplet, *godo.Response, error) {
					return []godo.Droplet{}, nil, errors.New("asdf")
				},
			},
			expectedError: errors.New("asdf"),
		},
		{
			name: "current page errors",
			ds: &stubDropletService{
				listTag: func(a string, b *godo.ListOptions) ([]godo.Droplet, *godo.Response, error) {
					resp := &godo.Response{}
					resp.Links = &godo.Links{Pages: &godo.Pages{Prev: "page=)", Last: "page="}}
					return []godo.Droplet{{Name: "foobar"}}, resp, nil
				},
			},
			expectedError: errors.New("parse page=): invalid URI for request"),
		},
	}

	for _, test := range tests {
		out, err := DropletListTags(test.ds, "access")
		if !reflect.DeepEqual(err, test.expectedError) {
			if err.Error() != test.expectedError.Error() {
				t.Logf("want:%v", test.expectedError)
				t.Logf("got:%v", err)
				t.Fatalf("test case failed: %s", test.name)
			}
		}
		if !reflect.DeepEqual(out, test.expectedDroplets) {
			t.Logf("want:%v", test.expectedDroplets)
			t.Logf("got:%v", out)
			t.Fatalf("test case failed: %s", test.name)
		}
	}
}

func TestSortDroplets(t *testing.T) {
	tests := []struct {
		name    string
		droplet godo.Droplet
		exp     map[string][]string
	}{
		{
			name: "no private iface",
			droplet: godo.Droplet{
				Region: &godo.Region{
					Slug: "nyc1",
				},
				Networks: &godo.Networks{
					V4: []godo.NetworkV4{
						godo.NetworkV4{IPAddress: "192.168.0.0", Type: "public"},
					},
				},
			},
			exp: map[string][]string{},
		},
		{
			name: "private iface",
			droplet: godo.Droplet{
				Region: &godo.Region{
					Slug: "nyc1",
				},
				Networks: &godo.Networks{
					V4: []godo.NetworkV4{
						godo.NetworkV4{IPAddress: "192.168.0.0", Type: "private"},
					},
				},
			},
			exp: map[string][]string{
				"nyc1": []string{"192.168.0.0"},
			},
		},
	}

	for _, test := range tests {
		out := SortDroplets([]godo.Droplet{test.droplet})
		if !reflect.DeepEqual(out, test.exp) {
			t.Logf("want:%v", test.exp)
			t.Logf("got:%v", out)
			t.Fatalf("test case failed: %s", test.name)
		}
	}
}

func TestPublicDroplets(t *testing.T) {
	tests := []struct {
		name    string
		droplet godo.Droplet
		exp     []string
	}{
		{
			name: "no public iface",
			droplet: godo.Droplet{
				Region: &godo.Region{
					Slug: "nyc1",
				},
				Networks: &godo.Networks{
					V4: []godo.NetworkV4{
						godo.NetworkV4{IPAddress: "192.168.0.0", Type: "private"},
					},
				},
			},
			exp: []string{},
		},
		{
			name: "public iface",
			droplet: godo.Droplet{
				Region: &godo.Region{
					Slug: "nyc1",
				},
				Networks: &godo.Networks{
					V4: []godo.NetworkV4{
						godo.NetworkV4{IPAddress: "192.168.0.0", Type: "public"},
					},
				},
			},
			exp: []string{"192.168.0.0"},
		},
	}

	for _, test := range tests {
		out := PublicDroplets([]godo.Droplet{test.droplet})
		if !reflect.DeepEqual(out, test.exp) {
			t.Logf("want:%v", test.exp)
			t.Logf("got:%v", out)
			t.Fatalf("test case failed: %s", test.name)
		}
	}
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
