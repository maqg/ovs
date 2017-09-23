package plugins

import (
	"octlink/ovs/utils"
	"octlink/ovs/utils/merrors"
	"octlink/ovs/utils/vyos"
)

// Dns for dns sturcture
type Dns struct {
	DnsAddress   string `json:"dnsAddress"`
	PublicNicMac string `json:"publicNicMac"`
}

// AddDns to add dns
func (d *Dns) AddDns() int {

	tree := vyos.NewParserFromShowConfiguration().Tree

	eth, err := utils.GetNicNameByMac(d.PublicNicMac)
	if err != nil {
		logger.Panicf("get nic name by mac %s error %s\n", d.PublicNicMac, err)
		return merrors.ErrBadParas
	}

	tree.Setf("service dns forwarding listen-on %s", eth)

	tree.Setf("service dns forwarding name-server %s", d.DnsAddress)

	tree.Apply(false)

	return 0
}

// GetAllDns by condition
func GetAllDns() []*Dns {

	tree := vyos.NewParserFromShowConfiguration().Tree

	dns := new(Dns)

	rule := tree.Getf("service dns forwarding")

	logger.Debugf("rule: %s\n", rule)

	return []*Dns{
		dns,
	}
}

// DeleteDns to delete dns
func (d *Dns) DeleteDns() int {

	tree := vyos.NewParserFromShowConfiguration().Tree

	tree.Deletef("service dns forwarding name-server %s", d.DnsAddress)

	tree.Apply(false)

	return merrors.ErrSuccess
}

/*

const (
	REMOVE_DNS_PATH = "/removedns"
	SET_DNS_PATH = "/setdns"
)

type dnsInfo struct {
	DnsAddress string `json:"dnsAddress"`
	NicMac string `json:"nicMac"`
}

type setDnsCmd struct {
	Dns []dnsInfo `json:"dns"`
}

type removeDnsCmd struct {
	Dns []dnsInfo `json:"dns"`
}


func makeDnsFirewallRuleDescription(nicname string) string {
	return fmt.Sprintf("DNS-for-%s", nicname)
}

func setDnsHandler(ctx *server.CommandContext) interface{} {
	tree := server.NewParserFromShowConfiguration().Tree

	cmd := &setDnsCmd{}
	ctx.GetCommand(cmd)

	dnsByMac := make(map[string][]dnsInfo)
	for _, info := range cmd.Dns {
		dns := dnsByMac[info.NicMac]
		if dns == nil {
			dns = make([]dnsInfo, 0)
		}
		dns = append(dns, info)
		dnsByMac[info.NicMac] = dns
	}

	for mac, dns := range dnsByMac {
		for _, info := range dns {
			tree.Setf("service dns forwarding name-server %s", info.DnsAddress)
		}
		eth, err := utils.GetNicNameByMac(mac); utils.PanicOnError(err)
		tree.Setf("service dns forwarding listen-on %s", eth)


		des := makeDnsFirewallRuleDescription(eth)
		if r := tree.FindFirewallRuleByDescription(eth, "local", des); r == nil {
			tree.SetFirewallOnInterface(eth, "local",
				fmt.Sprintf("description %v", des),
				"destination port 53",
				"protocol tcp_udp",
				"action accept",
			)

			tree.AttachFirewallToInterface(eth, "local")
		}
	}

	tree.Apply(false)

	return nil
}

func removeDnsHandler(ctx *server.CommandContext) interface{} {
	tree := server.NewParserFromShowConfiguration().Tree

	cmd := &removeDnsCmd{}
	ctx.GetCommand(cmd)

	for _, info := range cmd.Dns {
		tree.Deletef("service dns forwarding name-server %s", info.DnsAddress)
	}

	tree.Apply(false)
	return nil
}

func DnsEntryPoint() {
	server.RegisterAsyncCommandHandler(SET_DNS_PATH, server.VyosLock(setDnsHandler))
	server.RegisterAsyncCommandHandler(REMOVE_DNS_PATH, server.VyosLock(removeDnsHandler))
}

*/
