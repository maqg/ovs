package plugins

import (
	"fmt"
	"octlink/ovs/utils"
	"octlink/ovs/utils/merrors"
	"octlink/ovs/utils/vyos"
)

const (
	// VrConfigureNic Nic Vr Configure
	VrConfigureNic = "/configurenic"

	// VrRemoveNicPath VR remove nic path
	VrRemoveNicPath = "/removenic"

	// VrSSHPort for ssh port default
	VrSSHPort = 22

	// VrServicePort for vr service use
	VrServicePort = 3443
)

// IfInfo for Basic IfInfo Structure
type IfInfo struct {
	Name    string `json:"name"`
	IP      string `json:"ip"`
	Netmask string `json:"netmask"`
	Gateway string `json:"gateway"`
	Mac     string `json:"mac"`
}

// ConfigureNic by ifinfo
func (nic *IfInfo) ConfigureNic() int {

	tree := vyos.NewParserFromShowConfiguration().Tree

	nicname, err := utils.GetNicNameByMac(nic.Mac)
	utils.PanicOnError(err)
	cidr := utils.NetmaskToCIDR(nic.Netmask)
	utils.PanicOnError(err)

	addr := fmt.Sprintf("%v/%v", nic.IP, cidr)
	tree.SetfWithoutCheckExisting("interfaces ethernet %s address %v", nicname, addr)
	tree.SetfWithoutCheckExisting("interfaces ethernet %s duplex auto", nicname)
	tree.SetfWithoutCheckExisting("interfaces ethernet %s smp_affinity auto", nicname)
	tree.SetfWithoutCheckExisting("interfaces ethernet %s speed auto", nicname)

	tree.SetFirewallOnInterface(nicname, "local",
		"action accept",
		"state established enable",
		"state related enable",
		fmt.Sprintf("destination address %v", nic.IP),
	)
	tree.SetFirewallOnInterface(nicname, "local",
		"action accept",
		"protocol icmp",
		fmt.Sprintf("destination address %v", nic.IP),
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
		fmt.Sprintf("destination port %v", VrSSHPort),
		fmt.Sprintf("destination address %v", nic.IP),
		"protocol tcp",
		"action accept",
	)

	// config nic for api test
	tree.SetFirewallOnInterface(nicname, "local",
		fmt.Sprintf("destination port %v", VrServicePort),
		fmt.Sprintf("destination address %v", nic.IP),
		"protocol tcp",
		"action accept",
	)

	tree.SetFirewallDefaultAction(nicname, "local", "reject")
	tree.SetFirewallDefaultAction(nicname, "in", "reject")

	tree.AttachFirewallToInterface(nicname, "local")
	tree.AttachFirewallToInterface(nicname, "in")

	tree.Apply(false)

	return 0
}

// ConfigureNics for nic infos config
func ConfigureNics(nics []*IfInfo) int {
	for _, nic := range nics {
		nic.ConfigureNic()
	}
	return merrors.ErrSuccess
}

// RemoveNic by ifinfo
func (nic *IfInfo) RemoveNic() int {

	tree := vyos.NewParserFromShowConfiguration().Tree

	nicname, err := utils.GetNicNameByMac(nic.Mac)
	utils.PanicOnError(err)
	tree.Deletef("interfaces ethernet %s", nicname)
	tree.Deletef("firewall name %s.in", nicname)
	tree.Deletef("firewall name %s.local", nicname)

	tree.Apply(false)

	return merrors.ErrSuccess
}

// RemoveNics for nics removing
func RemoveNics(nics []*IfInfo) int {
	for _, nic := range nics {
		nic.RemoveNic()
	}
	return merrors.ErrSuccess
}

// GetNics by condition
func GetNics() []*IfInfo {

	ifs := make([]*IfInfo, 0)

	nics, err := utils.GetAllNics()
	if err != nil {
		fmt.Printf("get all nics error\n")
		return ifs
	}

	for _, nic := range nics {
		ifinfo := &IfInfo{
			Name: nic.Name,
			Mac:  nic.Mac,
		}
		ip, netmask, _, err := utils.GetNicInfo(nic.Name)
		if err == nil {
			ifinfo.IP = ip
			ifinfo.Netmask = netmask
		}
		ifs = append(ifs, ifinfo)
	}

	tree := vyos.NewParserFromShowConfiguration().Tree

	if rs := tree.Getf("interfaces ethernet"); rs != nil {
		logger.Debugf("ethernet config %s\n", rs.String())
	}

	return ifs
}
