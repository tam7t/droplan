package main

import (
	"errors"
	"testing"

	. "github.com/franela/goblin"
)

func TestTables(t *testing.T) {
	var sipt *stubIPTables

	g := Goblin(t)

	g.Describe(`Setup`, func() {
		g.BeforeEach(func() {
			sipt = newStubIPTables()
		})

		g.Describe(`when new chain returns an error`, func() {
			g.Describe(`because the chain already exists`, func() {
				g.BeforeEach(func() {
					sipt.newChain = func(a, b string) error {
						return errors.New("exit status 1: iptables: Chain already exists.\n")
					}
				})

				g.It(`does not error`, func() {
					g.Assert(Setup(sipt, `eth1`)).Equal(nil)
				})
			})

			g.Describe(`because of some other error`, func() {
				g.BeforeEach(func() {
					sipt.newChain = func(a, b string) error {
						return errors.New("something bad")
					}
				})

				g.It(`does not error`, func() {
					g.Assert(Setup(sipt, `eth1`)).Equal(errors.New("something bad"))
				})
			})
		})

		g.It(`creates a new chain`, func() {
			sipt.newChain = func(a, b string) error {
				g.Assert(a).Equal(`filter`)
				g.Assert(b).Equal(`dolan-peers`)
				return nil
			}
			Setup(sipt, `eth1`)
		})

		g.Describe(`when adding the chain to the interface errors`, func() {
			g.It(`returns the error`, func() {
				sipt.appendUnique = func(a, b string, c ...string) error {
					if a == `filter` && b == `INPUT` && len(c) == 4 {
						if c[0] == `-i` && c[1] == `eth1` && c[2] == `-j` && c[3] == `dolan-peers` {
							return errors.New(`bad add chain`)
						}
					}
					return nil
				}
				g.Assert(Setup(sipt, `eth1`)).Equal(errors.New(`bad add chain`))
			})
		})

		g.Describe(`when adding the deny rule errors`, func() {
			g.It(`returns the error`, func() {
				sipt.appendUnique = func(a, b string, c ...string) error {
					if a == `filter` && b == `INPUT` && len(c) == 4 {
						if c[0] == `-i` && c[1] == `eth1` && c[2] == `-j` && c[3] == `DROP` {
							return errors.New(`bad deny rule`)
						}
					}
					return nil
				}
				g.Assert(Setup(sipt, `eth1`)).Equal(errors.New(`bad deny rule`))
			})
		})

		g.It(`adds the dolan-peer chain and deny to the interface`, func() {
			var params [][]string

			sipt.appendUnique = func(a, b string, c ...string) error {
				args := []string{a, b}
				args = append(args, c...)
				params = append(params, args)
				return nil
			}

			Setup(sipt, `eth1`)

			g.Assert(params).Equal([][]string{
				[]string{`filter`, `INPUT`, `-i`, `eth1`, `-j`, `dolan-peers`},
				[]string{`filter`, `INPUT`, `-i`, `eth1`, `-j`, `DROP`},
			})
		})
	})

	g.Describe(`UpdatePeers`, func() {
		var peers []string

		g.BeforeEach(func() {
			sipt = newStubIPTables()
		})

		g.Describe(`with no peers`, func() {
			g.BeforeEach(func() {
				peers = []string{}
			})

			g.It(`clears the chain`, func() {
				sipt.clearChain = func(a, b string) error {
					g.Assert(a).Equal(`filter`)
					g.Assert(b).Equal(`dolan-peers`)
					return nil
				}
				UpdatePeers(sipt, peers)
			})

			g.It(`appends nothing to the chain`, func() {
				sipt.append = func(a, b string, c ...string) error {
					g.Fail(`append should not be called`)
					return nil
				}
				UpdatePeers(sipt, peers)
			})
		})

		g.Describe(`with 1 peer`, func() {
			g.BeforeEach(func() {
				peers = []string{`peer1`}
			})

			g.It(`clears the chain`, func() {
				sipt.clearChain = func(a, b string) error {
					g.Assert(a).Equal(`filter`)
					g.Assert(b).Equal(`dolan-peers`)
					return nil
				}
				UpdatePeers(sipt, peers)
			})

			g.It(`appends the peer to the chain`, func() {
				var params [][]string

				sipt.append = func(a, b string, c ...string) error {
					args := []string{a, b}
					args = append(args, c...)
					params = append(params, args)
					return nil
				}

				UpdatePeers(sipt, peers)

				g.Assert(params).Equal([][]string{
					[]string{`filter`, `dolan-peers`, `-s`, `peer1`, `-j`, `ACCEPT`},
				})
			})

			g.Describe(`when clearing the chain errors`, func() {
				g.It(`returns the error`, func() {
					sipt.clearChain = func(a, b string) error {
						return errors.New(`bad clear`)
					}
					g.Assert(UpdatePeers(sipt, peers)).Equal(errors.New(`bad clear`))
				})
			})

			g.Describe(`when appending errors`, func() {
				g.It(`returns the error`, func() {
					sipt.append = func(a, b string, c ...string) error {
						return errors.New(`bad append`)
					}
					g.Assert(UpdatePeers(sipt, peers)).Equal(errors.New(`bad append`))
				})
			})
		})

		g.Describe(`with many peers`, func() {
			g.BeforeEach(func() {
				peers = []string{`peer1`, `peer2`}
			})

			g.It(`clears the chain`, func() {
				sipt.clearChain = func(a, b string) error {
					g.Assert(a).Equal(`filter`)
					g.Assert(b).Equal(`dolan-peers`)
					return nil
				}
				UpdatePeers(sipt, peers)
			})

			g.It(`appends the peers to the chain`, func() {
				var params [][]string

				sipt.append = func(a, b string, c ...string) error {
					args := []string{a, b}
					args = append(args, c...)
					params = append(params, args)
					return nil
				}

				UpdatePeers(sipt, peers)

				g.Assert(params).Equal([][]string{
					[]string{`filter`, `dolan-peers`, `-s`, `peer1`, `-j`, `ACCEPT`},
					[]string{`filter`, `dolan-peers`, `-s`, `peer2`, `-j`, `ACCEPT`},
				})
			})
		})
	})
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
