package main

import (
	"errors"
	"reflect"
	"testing"
)

func TestSetup(t *testing.T) {
	count := 0

	tests := []struct {
		name string
		ipt  IPTables
		exp  error
	}{
		{
			name: "chain exists",
			ipt: &stubIPTables{
				newChain: func(a, b string) error {
					return errors.New("exit status 1: iptables: Chain already exists.\n")
				},
				clearChain:   func(string, string) error { return nil },
				appendUnique: func(string, string, ...string) error { return nil },
				append:       func(string, string, ...string) error { return nil },
			},
			exp: nil,
		},
		{
			name: "chain error",
			ipt: &stubIPTables{
				newChain: func(a, b string) error {
					return errors.New("something bad")
				},
				clearChain:   func(string, string) error { return nil },
				appendUnique: func(string, string, ...string) error { return nil },
				append:       func(string, string, ...string) error { return nil },
			},
			exp: errors.New("something bad"),
		},
		{
			name: "adds the droplan-peers filter",
			ipt: &stubIPTables{
				newChain: func(a, b string) error {
					if a == "filter" && b == "droplan-peers" {
						return nil
					}
					return errors.New("bad params")
				},
				clearChain:   func(string, string) error { return nil },
				appendUnique: func(string, string, ...string) error { return nil },
				append:       func(string, string, ...string) error { return nil },
			},
			exp: nil,
		},
		{
			name: "when adding the chain to the interface errors",
			ipt: &stubIPTables{
				newChain:   func(string, string) error { return nil },
				clearChain: func(string, string) error { return nil },
				appendUnique: func(a, b string, c ...string) error {
					if a == "filter" && b == "INPUT" && len(c) == 4 {
						if c[0] == "-i" && c[1] == "eth1" && c[2] == "-j" && c[3] == "droplan-peers" {
							return errors.New("bad add chain")
						}
					}
					return nil
				},
				append: func(string, string, ...string) error { return nil },
			},
			exp: errors.New("bad add chain"),
		},
		{
			name: "when adding the established connection rule errors",
			ipt: &stubIPTables{
				newChain:   func(string, string) error { return nil },
				clearChain: func(string, string) error { return nil },
				appendUnique: func(a, b string, c ...string) error {
					if a == "filter" && b == "INPUT" && len(c) == 8 {
						return errors.New("bad connect rule")
					}
					return nil
				},
				append: func(string, string, ...string) error { return nil },
			},
			exp: errors.New("bad connect rule"),
		},
		{
			name: "when adding the deny rule errors",
			ipt: &stubIPTables{
				newChain:   func(string, string) error { return nil },
				clearChain: func(string, string) error { return nil },
				appendUnique: func(a, b string, c ...string) error {
					if a == "filter" && b == "INPUT" && len(c) == 4 {
						if c[0] == "-i" && c[1] == "eth1" && c[2] == "-j" && c[3] == "DROP" {
							return errors.New("bad deny rule")
						}
					}
					return nil
				},
				append: func(string, string, ...string) error { return nil },
			},
			exp: errors.New("bad deny rule"),
		},
		{
			name: "append rules in order",
			ipt: &stubIPTables{
				newChain:   func(string, string) error { return nil },
				clearChain: func(string, string) error { return nil },
				appendUnique: func(a, b string, c ...string) error {
					if a == "filter" && b == "INPUT" && len(c) == 4 {
						if c[0] == "-i" && c[1] == "eth1" && c[2] == "-j" && c[3] == "DROP" {
							return errors.New("bad deny rule")
						}
					}
					return nil
				},
				append: func(string, string, ...string) error { return nil },
			},
			exp: errors.New("bad deny rule"),
		},
		{
			name: "adds peer chain and drop interface",
			ipt: &stubIPTables{
				newChain:   func(string, string) error { return nil },
				clearChain: func(string, string) error { return nil },
				appendUnique: func(a, b string, c ...string) error {
					defer func() { count++ }()
					switch count {
					case 0:
						if a == "filter" && b == "INPUT" && reflect.DeepEqual(c, []string{"-i", "eth1", "-j", "droplan-peers"}) {
							return nil
						} else {
							return errors.New("bad case 0")
						}
					case 1:
						if a == "filter" && b == "INPUT" && reflect.DeepEqual(c, []string{"-i", "eth1", "-m", "conntrack", "--ctstate", "ESTABLISHED,RELATED", "-j", "ACCEPT"}) {
							return nil
						} else {
							return errors.New("bad case 1")
						}
					case 2:
						if a == "filter" && b == "INPUT" && reflect.DeepEqual(c, []string{"-i", "eth1", "-j", "DROP"}) {
							return nil
						} else {
							return errors.New("bad case 2")
						}
					default:
						return errors.New("bad input")
					}
					return nil
				},
				append: func(string, string, ...string) error { return nil },
			},
		},
	}

	for _, test := range tests {
		out := Setup(test.ipt, "eth1")
		if !reflect.DeepEqual(out, test.exp) {
			t.Logf("want:%v", test.exp)
			t.Logf("got:%v", out)
			t.Fatalf("test case failed: %s", test.name)
		}
	}
}

func TestUpdatePeers(t *testing.T) {
	count := 0

	tests := []struct {
		name  string
		ipt   IPTables
		peers []string
		exp   error
	}{
		{
			name: "clears the chain",
			ipt: &stubIPTables{
				newChain: func(string, string) error { return nil },
				clearChain: func(a, b string) error {
					if a == "filter" && b == "droplan-peers" {
						return nil
					}
					return errors.New("bad clear chain args")
				},
				appendUnique: func(string, string, ...string) error { return nil },
				append:       func(string, string, ...string) error { return nil },
			},
		},
		{
			name: "does not append anythign if peers are empty",
			ipt: &stubIPTables{
				newChain:     func(string, string) error { return nil },
				clearChain:   func(string, string) error { return nil },
				appendUnique: func(string, string, ...string) error { return nil },
				append: func(string, string, ...string) error {
					return errors.New("no peers should be appended")
				},
			},
		},
		{
			name: "when clearing the chain errors",
			ipt: &stubIPTables{
				newChain: func(string, string) error { return nil },
				clearChain: func(string, string) error {
					return errors.New("clear chain error")
				},
				appendUnique: func(string, string, ...string) error { return nil },
				append:       func(string, string, ...string) error { return nil },
			},
			exp: errors.New("clear chain error"),
		},
		{
			name: "when appending to the chain errors",
			ipt: &stubIPTables{
				newChain:     func(string, string) error { return nil },
				clearChain:   func(string, string) error { return nil },
				appendUnique: func(string, string, ...string) error { return nil },
				append: func(string, string, ...string) error {
					return errors.New("peer append error")
				},
			},
			peers: []string{"peer1"},
			exp:   errors.New("peer append error"),
		},
		{
			name: "adds peer chain and drop interface",
			ipt: &stubIPTables{
				newChain:     func(string, string) error { return nil },
				clearChain:   func(string, string) error { return nil },
				appendUnique: func(string, string, ...string) error { return nil },
				append: func(a, b string, c ...string) error {
					defer func() { count++ }()
					switch count {
					case 0:
						if a == "filter" && b == "droplan-peers" && reflect.DeepEqual(c, []string{"-s", "peer1", "-j", "ACCEPT"}) {
							return nil
						} else {
							return errors.New("bad case 0")
						}
					case 1:
						if a == "filter" && b == "droplan-peers" && reflect.DeepEqual(c, []string{"-s", "peer2", "-j", "ACCEPT"}) {
							return nil
						} else {
							return errors.New("bad case 1")
						}
					case 2:
						if a == "filter" && b == "droplan-peers" && reflect.DeepEqual(c, []string{"-s", "peer3", "-j", "ACCEPT"}) {
							return nil
						} else {
							return errors.New("bad case 2")
						}
					default:
						return errors.New("bad input")
					}
					return nil
				},
			},
			peers: []string{"peer1", "peer2", "peer3"},
		},
	}

	for _, test := range tests {
		out := UpdatePeers(test.ipt, test.peers)
		if !reflect.DeepEqual(out, test.exp) {
			t.Logf("want:%v", test.exp)
			t.Logf("got:%v", out)
			t.Fatalf("test case failed: %s", test.name)
		}
	}
}

func newStubIPTables() *stubIPTables {
	return &stubIPTables{
		newChain:     func(string, string) error { return nil },
		clearChain:   func(string, string) error { return nil },
		appendUnique: func(string, string, ...string) error { return nil },
		append:       func(string, string, ...string) error { return nil },
	}
}

type stubIPTables struct {
	newChain     func(string, string) error
	clearChain   func(string, string) error
	appendUnique func(string, string, ...string) error
	append       func(string, string, ...string) error
}

func (sipt *stubIPTables) ClearChain(a, b string) error {
	return sipt.clearChain(a, b)
}

func (sipt *stubIPTables) Append(a, b string, c ...string) error {
	return sipt.append(a, b, c...)
}

func (sipt *stubIPTables) AppendUnique(a, b string, c ...string) error {
	return sipt.appendUnique(a, b, c...)
}

func (sipt *stubIPTables) NewChain(a, b string) error {
	return sipt.newChain(a, b)
}
