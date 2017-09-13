package plugins

import (
	"fmt"
	"octlink/ovs/utils"
	"octlink/ovs/utils/vyos"
	"strings"
)

// EipInfo base structure
type EipInfo struct {
	VipIP              string `json:"vipIp"`
	PrivateMac         string `json:"privateMac"`
	GuestIP            string `json:"guestIp"`
	PublicMac          string `json:"publicMac"`
	SnatInboundTraffic bool   `json:"snatInboundTraffic"`
}

func makeEipDescription(info *EipInfo) string {
	return fmt.Sprintf("EIP-%v-%v-%v", info.VipIP, info.GuestIP, info.PrivateMac)
}

func makeEipDescriptionForPrivateMac(info *EipInfo) string {
	return fmt.Sprintf("EIP-%v-%v-%v-private", info.VipIP, info.GuestIP, info.PrivateMac)
}

func setEip(tree *vyos.ConfigTree, eip *EipInfo) {
	des := makeEipDescription(eip)
	priDes := makeEipDescriptionForPrivateMac(eip)
	nicname, err := utils.GetNicNameByIP(eip.VipIP)
	if err != nil && eip.PublicMac != "" {
		nicname, err = utils.GetNicNameByMac(eip.PublicMac)
	}
	utils.PanicOnError(err)

	prinicname, err := utils.GetNicNameByMac(eip.PrivateMac)
	utils.PanicOnError(err)

	if r := tree.FindSnatRuleDescription(des); r == nil {
		tree.SetSnat(
			fmt.Sprintf("description %v", des),
			fmt.Sprintf("outbound-interface %v", nicname),
			fmt.Sprintf("source address %v", eip.GuestIP),
			fmt.Sprintf("translation address %v", eip.VipIP),
		)
	}

	if r := tree.FindSnatRuleDescription(priDes); r == nil {
		tree.SetSnat(
			fmt.Sprintf("description %v", priDes),
			fmt.Sprintf("outbound-interface %v", prinicname),
			fmt.Sprintf("source address %v", eip.GuestIP),
			fmt.Sprintf("translation address %v", eip.VipIP),
		)
	}

	if r := tree.FindDnatRuleDescription(des); r == nil {
		tree.SetDnat(
			fmt.Sprintf("description %v", des),
			fmt.Sprintf("inbound-interface any"),
			fmt.Sprintf("destination address %v", eip.VipIP),
			fmt.Sprintf("translation address %v", eip.GuestIP),
		)
	}

	if r := tree.FindFirewallRuleByDescription(nicname, "in", des); r == nil {
		tree.SetFirewallOnInterface(nicname, "in",
			fmt.Sprintf("description %v", des),
			fmt.Sprintf("destination address %v", eip.GuestIP),
			"state new enable",
			"state established enable",
			"state related enable",
			"action accept",
		)

		tree.AttachFirewallToInterface(nicname, "in")
	}

	if r := tree.FindFirewallRuleByDescription(prinicname, "in", des); r == nil {
		tree.SetFirewallOnInterface(prinicname, "in",
			fmt.Sprintf("description %v", des),
			fmt.Sprintf("source address %v", eip.GuestIP),
			"state new enable",
			"state established enable",
			"state related enable",
			"action accept",
		)

		tree.AttachFirewallToInterface(prinicname, "in")
	}
}

func deleteEip(tree *vyos.ConfigTree, eip *EipInfo) {
	des := makeEipDescription(eip)
	priDes := makeEipDescriptionForPrivateMac(eip)
	nicname, err := utils.GetNicNameByIP(eip.VipIP)
	if err != nil && eip.PublicMac != "" {
		nicname, err = utils.GetNicNameByMac(eip.PublicMac)
	}
	utils.PanicOnError(err)

	if r := tree.FindSnatRuleDescription(des); r != nil {
		r.Delete()
	}

	if r := tree.FindSnatRuleDescription(priDes); r != nil {
		r.Delete()
	}

	if r := tree.FindDnatRuleDescription(des); r != nil {
		r.Delete()
	}

	if r := tree.FindFirewallRuleByDescription(nicname, "in", des); r != nil {
		r.Delete()
	}

	prinicname, err := utils.GetNicNameByMac(eip.PrivateMac)
	utils.PanicOnError(err)
	if r := tree.FindFirewallRuleByDescription(prinicname, "in", des); r != nil {
		r.Delete()
	}
}

// CreateEip to remove eip
func CreateEip() int {
	eip := &EipInfo{}
	tree := vyos.NewParserFromShowConfiguration().Tree
	setEip(tree, eip)
	tree.Apply(false)

	return 0
}

func removeEip() interface{} {
	tree := vyos.NewParserFromShowConfiguration().Tree
	deleteEip(tree, &EipInfo{})
	tree.Apply(false)

	return nil
}

func syncEips(eips []*EipInfo) int {

	tree := vyos.NewParserFromShowConfiguration().Tree

	// delete all EIP related rules
	if rs := tree.Get("nat destination rule"); rs != nil {
		for _, r := range rs.Children() {
			if d := r.Get("description"); d != nil && strings.HasPrefix(d.Value(), "EIP") {
				r.Delete()
			}
		}
	}

	if rs := tree.Getf("nat source rule"); rs != nil {
		for _, r := range rs.Children() {
			if d := r.Get("description"); d != nil && strings.HasPrefix(d.Value(), "EIP") {
				r.Delete()
			}
		}
	}

	if rs := tree.Getf("firewall name"); rs != nil {
		for _, r := range rs.Children() {
			if rss := r.Get("rule"); rss != nil {
				for _, rr := range rss.Children() {
					if d := rr.Get("description"); d != nil && strings.HasPrefix(d.Value(), "EIP") {
						rr.Delete()
					}
				}
			}
		}
	}

	for _, eip := range eips {
		setEip(tree, eip)
	}

	tree.Apply(false)

	return 0
}
