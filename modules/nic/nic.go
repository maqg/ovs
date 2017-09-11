package plugin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"ovs/utils"

	log "github.com/Sirupsen/logrus"
)

const (
	VR_CONFIGURE_NIC     = "/configurenic"
	VR_REMOVE_NIC_PATH   = "/removenic"
	BOOTSTRAP_INFO_CACHE = "/home/vyos/zvr/bootstrap-info.json"
	DEFAULT_SSH_PORT     = 22
)

type nicInfo struct {
	Ip      string `json:"ip"`
	Netmask string `json:"netmask"`
	Gateway string `json:"gateway"`
	Mac     string `json:"Mac"`
}

type configureNicCmd struct {
	Nics []nicInfo `json:"nics"`
}

var bootstrapInfo map[string]interface{} = make(map[string]interface{})

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

func getSshPortFromBootInfo() float64 {
	content, err := ioutil.ReadFile(BOOTSTRAP_INFO_CACHE)
	utils.PanicOnError(err)
	if len(content) == 0 {
		log.Debugf("no content in %s, use default ssh port %d", BOOTSTRAP_INFO_CACHE, DEFAULT_SSH_PORT)
		return DEFAULT_SSH_PORT
	}

	if err := json.Unmarshal(content, &bootstrapInfo); err != nil {
		log.Debugf("can not parse info from %s, use default ssh port %d", BOOTSTRAP_INFO_CACHE, DEFAULT_SSH_PORT)
		return DEFAULT_SSH_PORT
	}

	return bootstrapInfo["sshPort"].(float64)
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

func ConfigureNicEntryPoint() {
	//server.RegisterAsyncCommandHandler(VR_CONFIGURE_NIC, server.VyosLock(configureNic))
	//server.RegisterAsyncCommandHandler(VR_REMOVE_NIC_PATH, server.VyosLock(removeNic))
}
