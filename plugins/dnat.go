package plugins

import (
	"fmt"
	"octlink/ovs/utils"
	"octlink/ovs/utils/merrors"
	"octlink/ovs/utils/vyos"
	"strings"
)

// Dnat for dnat sturcture
type Dnat struct {
	VipPortStart     int    `json:"vipPortStart"`
	VipPortEnd       int    `json:"vipPortEnd"`
	PrivatePortStart int    `json:"privatePortStart"`
	PrivatePortEnd   int    `json:"privatePortEnd"`
	ProtocolType     string `json:"protocolType"`
	VipIp            string `json:"vipIp"`
	PrivateIp        string `json:"privateIp"`
	PrivateNicMac    string `json:"privateNicMac"`
	AllowedCidr      string `json:"allowedCidr"`
}

func makeDnatDescription(dnat *Dnat) string {
	return fmt.Sprintf("%v-%v-%v-%v-%v-%v-%v", dnat.VipIp, dnat.VipPortStart, dnat.VipPortEnd, dnat.PrivateNicMac, dnat.PrivatePortStart, dnat.PrivatePortEnd, dnat.ProtocolType)
}

func setDnat(tree *vyos.ConfigTree, dnat *Dnat) {

	var sport string
	if dnat.VipPortStart == dnat.VipPortEnd {
		sport = fmt.Sprintf("%v", dnat.VipPortStart)
	} else {
		sport = fmt.Sprintf("%v-%v", dnat.VipPortStart, dnat.VipPortEnd)
	}
	var dport string
	if dnat.PrivatePortStart == dnat.PrivatePortEnd {
		dport = fmt.Sprintf("%v", dnat.PrivatePortStart)
	} else {
		dport = fmt.Sprintf("%v-%v", dnat.PrivatePortStart, dnat.PrivatePortEnd)
	}

	pubNicName, err := utils.GetNicNameByIP(dnat.VipIp)
	utils.PanicOnError(err)

	des := makeDnatDescription(dnat)
	if r := tree.FindDnatRuleDescription(des); r == nil {
		tree.SetDnat(
			fmt.Sprintf("description %v", des),
			fmt.Sprintf("destination address %v", dnat.VipIp),
			fmt.Sprintf("destination port %v", sport),
			fmt.Sprintf("inbound-interface any"),
			fmt.Sprintf("protocol %v", strings.ToLower(dnat.ProtocolType)),
			fmt.Sprintf("translation address %v", dnat.PrivateIp),
			fmt.Sprintf("translation port %v", dport),
		)
	}

	if fr := tree.FindFirewallRuleByDescription(pubNicName, "in", des); fr == nil {
		if dnat.AllowedCidr != "" && dnat.AllowedCidr != "0.0.0.0/0" {
			tree.SetFirewallOnInterface(pubNicName, "in",
				"action reject",
				fmt.Sprintf("source address !%v", dnat.AllowedCidr),
				fmt.Sprintf("description %v", des),
				// NOTE: the destination is private IP
				// because the destination address is changed by the dnat rule
				fmt.Sprintf("destination address %v", dnat.PrivateIp),
				fmt.Sprintf("destination port %v", dport),
				fmt.Sprintf("protocol %s", strings.ToLower(dnat.ProtocolType)),
				"state new enable",
			)
		} else {
			tree.SetFirewallOnInterface(pubNicName, "in",
				"action accept",
				fmt.Sprintf("description %v", des),
				fmt.Sprintf("destination address %v", dnat.PrivateIp),
				fmt.Sprintf("destination port %v", dport),
				fmt.Sprintf("protocol %s", strings.ToLower(dnat.ProtocolType)),
				"state new enable",
			)
		}
	}

	tree.AttachFirewallToInterface(pubNicName, "in")

}

// AddDnat for add dnat
func (dnat *Dnat) AddDnat() int {

	tree := vyos.NewParserFromShowConfiguration().Tree
	setDnat(tree, dnat)
	tree.Apply(false)

	return 0
}

func deleteDnat(tree *vyos.ConfigTree, dnat *Dnat) {

	des := makeDnatDescription(dnat)
	if r := tree.FindDnatRuleDescription(des); r != nil {
		r.Delete()
	}

	pubNicName, err := utils.GetNicNameByIP(dnat.VipIp)
	utils.PanicOnError(err)

	if fr := tree.FindFirewallRuleByDescription(pubNicName, "in", des); fr != nil {
		fr.Delete()
	}
}

// RemoveDnat for remove dnat
func (dnat *Dnat) RemoveDnat() int {

	tree := vyos.NewParserFromShowConfiguration().Tree

	deleteDnat(tree, dnat)
	tree.Apply(false)

	return 0
}

// RemoveDnats to remove eips from VR
func RemoveDnats(dnats []*Dnat) int {

	tree := vyos.NewParserFromShowConfiguration().Tree

	for _, dnat := range dnats {
		deleteDnat(tree, dnat)
	}

	tree.Apply(false)

	return merrors.ErrSuccess
}

// SyncDnats to sync all eips
func SyncDnats(dnats []*Dnat) int {

	return 0
}

// GetAllDnats get all dnats config
func GetAllDnats() []*Dnat {

	var dnats []*Dnat
	tree := vyos.NewParserFromShowConfiguration().Tree

	if rs := tree.Get("nat destination rule"); rs != nil {
		for _, r := range rs.Children() {
			if d := r.Get("description"); d != nil && (strings.HasSuffix(d.Value(), "TCP") || strings.HasSuffix(d.Value(), "UDP")) {

				dnat := new(Dnat)

				descList := strings.Split(d.Value(), "-")
				if len(descList) != 7 {
					continue
				}

				dnat.VipIp = descList[0]
				dnat.VipPortStart = utils.StringToInt(descList[1])
				dnat.VipPortEnd = utils.StringToInt(descList[2])
				dnat.PrivateNicMac = descList[3]
				dnat.PrivatePortStart = utils.StringToInt(descList[4])
				dnat.PrivatePortEnd = utils.StringToInt(descList[5])
				dnat.ProtocolType = descList[6]

				dnat.PrivateIp = r.Get("translation address").Value()

				dnats = append(dnats, dnat)
			}
		}
	}

	return dnats
}

// GetDnat to get dnat by privateMac
func GetDnat(privateNicMac string) (*Dnat, int) {

	dnats := GetAllDnats()
	for _, dnat := range dnats {
		if dnat.PrivateNicMac == privateNicMac {
			return dnat, merrors.ErrSuccess
		}
	}
	return nil, merrors.ErrSegmentNotExist
}

/*

const (
	CREATE_PORT_FORWARDING_PATH = "/createportforwarding"
	REVOKE_PORT_FORWARDING_PATH = "/revokeportforwarding"
	SYNC_PORT_FORWARDING_PATH = "/syncportforwarding"
)

type dnatInfo struct {
	VipPortStart int `json:"vipPortStart"`
	VipPortEnd int `json:"vipPortEnd"`
	PrivatePortStart int `json:"privatePortStart"`
	PrivatePortEnd int `json:"privatePortEnd"`
	ProtocolType string `json:"protocolType"`
	VipIp string `json:"vipIp"`
	PrivateIp string `json:"privateIp"`
	PrivateMac string `json:"privateMac"`
	AllowedCidr string `json:"allowedCidr"`
	SnatInboundTraffic bool `json:"snatInboundTraffic"`
}

type setDnatCmd struct {
	Rules []dnatInfo `json:"rules"`
}

type removeDnatCmd struct {
	Rules []dnatInfo `json:"rules"`
}

type syncDnatCmd struct {
	Rules []dnatInfo `json:"rules"`
}

func syncDnatHandler(ctx *server.CommandContext) interface{} {
	cmd := &syncDnatCmd{}
	ctx.GetCommand(cmd)

	tree := server.NewParserFromShowConfiguration().Tree
	tree.Delete("nat destination")
	setRuleInTree(tree, cmd.Rules)
	tree.Apply(false)
	return nil
}

func getRule(tree *server.VyosConfigTree, description string) *server.VyosConfigNode {
	rs := tree.Get("nat destination rule")
	if rs == nil {
		return nil
	}

	for _, r := range rs.Children() {
		if des := r.Get("description"); des != nil && des.Value() == description {
			return r
		}
	}

	return nil
}

func makeDnatDescription(r dnatInfo) string {
	return fmt.Sprintf("%v-%v-%v-%v-%v-%v-%v", r.VipIp, r.VipPortStart, r.VipPortEnd, r.PrivateMac, r.PrivatePortStart, r.PrivatePortEnd, r.ProtocolType)
}

func setRuleInTree(tree *server.VyosConfigTree, rules []dnatInfo) {
	for _, r := range rules {
		des := makeDnatDescription(r)
		if currentRule := getRule(tree, des); currentRule != nil {
			log.Debugf("dnat rule %s exists, skip it", des)
			continue
		}

		var sport string
		if r.VipPortStart == r.VipPortEnd {
			sport = fmt.Sprintf("%v", r.VipPortStart)
		} else {
			sport = fmt.Sprintf("%v-%v", r.VipPortStart, r.VipPortEnd)
		}
		var dport string
		if r.PrivatePortStart == r.PrivatePortEnd {
			dport = fmt.Sprintf("%v", r.PrivatePortStart)
		} else {
			dport = fmt.Sprintf("%v-%v", r.PrivatePortStart, r.PrivatePortEnd)
		}

		pubNicName, err := utils.GetNicNameByIp(r.VipIp); utils.PanicOnError(err)

		tree.SetDnat(
			fmt.Sprintf("description %v", des),
			fmt.Sprintf("destination address %v", r.VipIp),
			fmt.Sprintf("destination port %v", sport),
			fmt.Sprintf("inbound-interface any"),
			fmt.Sprintf("protocol %v", strings.ToLower(r.ProtocolType)),
			fmt.Sprintf("translation address %v", r.PrivateIp),
			fmt.Sprintf("translation port %v", dport),
		)

		if fr := tree.FindFirewallRuleByDescription(pubNicName, "in", des); fr == nil {
			if r.AllowedCidr != "" && r.AllowedCidr != "0.0.0.0/0" {
				tree.SetFirewallOnInterface(pubNicName, "in",
					"action reject",
					fmt.Sprintf("source address !%v", r.AllowedCidr),
					fmt.Sprintf("description %v", des),
					// NOTE: the destination is private IP
					// because the destination address is changed by the dnat rule
					fmt.Sprintf("destination address %v", r.PrivateIp),
					fmt.Sprintf("destination port %v", dport),
					fmt.Sprintf("protocol %s", strings.ToLower(r.ProtocolType)),
					"state new enable",
				)
			} else {
				tree.SetFirewallOnInterface(pubNicName, "in",
					"action accept",
					fmt.Sprintf("description %v", des),
					fmt.Sprintf("destination address %v", r.PrivateIp),
					fmt.Sprintf("destination port %v", dport),
					fmt.Sprintf("protocol %s", strings.ToLower(r.ProtocolType)),
					"state new enable",
				)
			}
		}

		tree.AttachFirewallToInterface(pubNicName, "in")
	}
}

func setDnatHandler(ctx *server.CommandContext) interface{} {
	cmd := &setDnatCmd{}
	ctx.GetCommand(cmd)

	tree := server.NewParserFromShowConfiguration().Tree
	setRuleInTree(tree, cmd.Rules)
	tree.Apply(false)

	return nil
}

func removeDnatHandler(ctx *server.CommandContext) interface{} {
	cmd := &removeDnatCmd{}
	ctx.GetCommand(cmd)

	tree := server.NewParserFromShowConfiguration().Tree
	for _, r := range cmd.Rules {
		des := makeDnatDescription(r)
		if c := getRule(tree, des); c != nil {
			c.Delete()
		}

		pubNicName, err := utils.GetNicNameByIp(r.VipIp); utils.PanicOnError(err)
		if fr := tree.FindFirewallRuleByDescription(pubNicName, "in", des); fr != nil {
			fr.Delete()
		}
	}
	tree.Apply(false)

	return nil
}

func DnatEntryPoint() {
	server.RegisterAsyncCommandHandler(CREATE_PORT_FORWARDING_PATH, server.VyosLock(setDnatHandler))
	server.RegisterAsyncCommandHandler(REVOKE_PORT_FORWARDING_PATH, server.VyosLock(removeDnatHandler))
	server.RegisterAsyncCommandHandler(SYNC_PORT_FORWARDING_PATH, server.VyosLock(syncDnatHandler))
}
*/
