package main

// IPTables interface for interacting with an iptables library. Declare it this
// way so that it is easy to dependency inject a mock.
type IPTables interface {
	ClearChain(string, string) error
	Append(string, string, ...string) error
	AppendUnique(string, string, ...string) error
	NewChain(string, string) error
}

// Setup creates a new iptables chain for holding peers and adds the chain and
// deny rules to the specified interface
func Setup(ipt IPTables, ipFace, chain string) error {
	var err error

	err = ipt.NewChain("filter", chain)
	if err != nil {
		if err.Error() != "exit status 1: iptables: Chain already exists.\n" {
			return err
		}
	}

	err = ipt.AppendUnique("filter", "INPUT", "-i", ipFace, "-j", chain)
	if err != nil {
		return err
	}
	// Do not drop connections when the `droplan-peers` chain is being updated
	err = ipt.AppendUnique("filter", "INPUT", "-i", ipFace, "-m", "conntrack", "--ctstate", "ESTABLISHED,RELATED", "-j", "ACCEPT")
	if err != nil {
		return err
	}
	err = ipt.AppendUnique("filter", "INPUT", "-i", ipFace, "-j", "DROP")
	if err != nil {
		return err
	}
	return nil
}

// UpdatePeers updates the droplan-peers chain in iptables with the specified
// peers
func UpdatePeers(ipt IPTables, peers []string, chain string) error {
	err := ipt.ClearChain("filter", chain)
	if err != nil {
		return err
	}

	for _, peer := range peers {
		err := ipt.Append("filter", chain, "-s", peer, "-j", "ACCEPT")
		if err != nil {
			return err
		}
	}
	return nil
}
