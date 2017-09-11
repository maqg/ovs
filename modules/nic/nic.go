package nic

import (
	"fmt"
	"octlink/ovs/utils"
	"octlink/ovs/utils/vyos"
)

const (
	// VrConfigureNic Nic Vr Configure
	VrConfigureNic = "/configurenic"

	// VrRemoveNicPath VR remove nic path
	VrRemoveNicPath = "/removenic"
)

// Nic for Basic Nic Structure
type Nic struct {
	IP      string `json:"ip"`
	Netmask string `json:"netmask"`
	Gateway string `json:"gateway"`
	Mac     string `json:"Mac"`
}

type configureNicCmd struct {
	Nics []Nic `json:"nics"`
}

func configureNic(ctx *vyos.CommandContext) interface{} {
	cmd := &configureNicCmd{}

	ctx.GetCommand(cmd)

	tree := vyos.NewParserFromShowConfiguration().Tree
	for _, nic := range cmd.Nics {
		nicname, err := utils.GetNicNameByMac(nic.Mac)
		utils.PanicOnError(err)
		cidr, err := utils.NetmaskToCIDR(nic.Netmask)
		utils.PanicOnError(err)
		addr := fmt.Sprintf("%v/%v", nic.Ip, cidr)
		tree.SetfWithoutCheckExisting("interfaces ethernet %s address %v", nicname, addr)
		tree.SetfWithoutCheckExisting("interfaces ethernet %s duplex auto", nicname)
		tree.SetfWithoutCheckExisting("interfaces ethernet %s smp_affinity auto", nicname)
		tree.SetfWithoutCheckExisting("interfaces ethernet %s speed auto", nicname)

		tree.SetFirewallOnInterface(nicname, "local",
			"action accept",
			"state established enable",
			"state related enable",
			fmt.Sprintf("destination address %v", nic.Ip),
		)
		tree.SetFirewallOnInterface(nicname, "local",
			"action accept",
			"protocol icmp",
			fmt.Sprintf("destination address %v", nic.Ip),
		)

		tree.SetFirewallOnInterface(nicname, "in",
			"action accept",
			"state established enable",
			"state related enable",
			"state new enable",
		)
		tree.SetFirewallOnInterface(nicname, "in",
			"action accept",
			"protocol icmp",
		)

		tree.SetFirewallOnInterface(nicname, "local",
			fmt.Sprintf("destination port %v", int(getSshPortFromBootInfo())),
			fmt.Sprintf("destination address %v", nic.Ip),
			"protocol tcp",
			"action accept",
		)

		tree.SetFirewallDefaultAction(nicname, "local", "reject")
		tree.SetFirewallDefaultAction(nicname, "in", "reject")

		tree.AttachFirewallToInterface(nicname, "local")
		tree.AttachFirewallToInterface(nicname, "in")
	}

	tree.Apply(false)
	return nil
}

func removeNic(ctx *vyos.CommandContext) interface{} {
	cmd := &configureNicCmd{}
	ctx.GetCommand(cmd)

	tree := vyos.NewParserFromShowConfiguration().Tree
	for _, nic := range cmd.Nics {
		nicname, err := utils.GetNicNameByMac(nic.Mac)
		utils.PanicOnError(err)
		tree.Deletef("interfaces ethernet %s", nicname)
		tree.Deletef("firewall name %s.in", nicname)
		tree.Deletef("firewall name %s.local", nicname)
	}
	tree.Apply(false)

	return nil
}

// GetInterfaces by condition
func GetInterfaces() []*Nic {
	return make([]*Nic, 0)
}
