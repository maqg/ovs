package main

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"octlink/ovs/utils"
	"octlink/ovs/utils/octlog"
	"octlink/ovs/utils/vyos"
	"os"
	"strings"
	"time"
)

const (
	VIRTIO_PORT_PATH     = "/dev/virtio-ports/applianceVm.vport"
	BOOTSTRAP_INFO_CACHE = "/home/vyos/rvm/bootstrap-info.json"
	TMP_LOCATION_FOR_ESX = "/tmp/bootstrap-info.json"
	// use this rule number to set a rule which confirm route entry work issue ZSTAC-6170
	ROUTE_STATE_NEW_ENABLE_FIREWALL_RULE_NUMBER = 9999
)

type nic struct {
	mac            string
	ip             string
	name           string
	netmask        string
	isDefaultRoute bool
	gateway        string
}

var bootstrapInfo map[string]interface{} = make(map[string]interface{})
var nics map[string]*nic = make(map[string]*nic)

func waitIptablesServiceOnline() {
	bash := utils.Bash{
		Command: "/sbin/iptables-save",
	}

	utils.LoopRunUntilSuccessOrTimeout(func() bool {
		err := bash.Run()
		if err != nil {
			octlog.Warn("iptables service seems not ready, %v", err)
		}
		return err == nil
	}, time.Duration(120)*time.Second, time.Duration(500)*time.Millisecond)
}

func waitVirtioPortOnline() {
	utils.LoopRunUntilSuccessOrTimeout(func() bool {
		ok, err := utils.PathExists(VIRTIO_PORT_PATH)
		utils.PanicOnError(err)
		if !ok {
			octlog.Warn("%s doesn't not exist, wait it ...", VIRTIO_PORT_PATH)
		}
		return ok
	}, time.Duration(120)*time.Second, time.Duration(500)*time.Millisecond)
}

func parseKvmBootInfo() {
	utils.LoopRunUntilSuccessOrTimeout(func() bool {
		content, err := ioutil.ReadFile(VIRTIO_PORT_PATH)
		utils.PanicOnError(err)
		if len(content) == 0 {
			octlog.Warn("no content in %s, it may not be ready, wait it ...", VIRTIO_PORT_PATH)
			return false
		}

		if err := json.Unmarshal(content, &bootstrapInfo); err != nil {
			panic(errors.Wrap(err, fmt.Sprintf("unable to JSON parse:\n %s", string(content))))
		}

		err = utils.MkdirForFile(BOOTSTRAP_INFO_CACHE, 0666)
		utils.PanicOnError(err)
		err = ioutil.WriteFile(BOOTSTRAP_INFO_CACHE, content, 0666)
		utils.PanicOnError(err)
		err = os.Chmod(BOOTSTRAP_INFO_CACHE, 0777)
		utils.PanicOnError(err)
		octlog.Debug("recieved bootstrap info:\n%s", string(content))
		return true
	}, time.Duration(300)*time.Second, time.Duration(1)*time.Second)
}

func resetVyos() {
	// clear all configuration in case someone runs 'save' command manually before,
	// to keep the vyos must be stateless

	// delete all interfaces
	tree := vyos.NewParserFromShowConfiguration().Tree
	tree.Delete("interfaces ethernet")
	tree.Apply(true)

	// reload default configuration
	vyos.RunVyosScriptAsUserVyos("load /opt/vyatta/etc/config.boot.default\nsave")
}

func configureVyos() {
	resetVyos()

	mgmtNic := bootstrapInfo["managementNic"].(map[string]interface{})
	if mgmtNic == nil {
		panic(errors.New("no field 'managementNic' in bootstrap info"))
	}

	eth0 := &nic{name: "eth0"}
	var ok bool
	eth0.mac, ok = mgmtNic["mac"].(string)
	utils.PanicIfError(ok, errors.New("cannot find 'mac' field for the management nic"))
	eth0.netmask, ok = mgmtNic["netmask"].(string)
	utils.PanicIfError(ok, errors.New("cannot find 'netmask' field for the management nic"))
	eth0.ip, ok = mgmtNic["ip"].(string)
	utils.PanicIfError(ok, errors.New("cannot find 'ip' field for the management nic"))
	eth0.isDefaultRoute = mgmtNic["isDefaultRoute"].(bool)
	eth0.gateway = mgmtNic["gateway"].(string)
	nics[eth0.name] = eth0

	otherNics := bootstrapInfo["additionalNics"].([]interface{})
	if otherNics != nil {
		for _, o := range otherNics {
			onic := o.(map[string]interface{})
			n := &nic{}
			n.name, ok = onic["deviceName"].(string)
			utils.PanicIfError(ok, fmt.Errorf("cannot find 'deviceName' field for the nic"))
			n.mac, ok = onic["mac"].(string)
			utils.PanicIfError(ok, errors.New("cannot find 'mac' field for the nic"))
			n.netmask, ok = onic["netmask"].(string)
			utils.PanicIfError(ok, fmt.Errorf("cannot find 'netmask' field for the nic[name:%s]", n.name))
			n.ip, ok = onic["ip"].(string)
			utils.PanicIfError(ok, fmt.Errorf("cannot find 'ip' field for the nic[name:%s]", n.name))
			n.gateway = onic["gateway"].(string)
			n.isDefaultRoute = onic["isDefaultRoute"].(bool)
			nics[n.name] = n
		}
	}

	type deviceName struct {
		expected string
		actual   string
		swap     string
	}

	devNames := make([]*deviceName, 0)

	// check integrity of nics
	for _, nic := range nics {
		utils.Assertf(nic.name != "", "name cannot be empty[mac:%s]", nic.mac)
		utils.Assertf(nic.ip != "", "ip cannot be empty[nicname: %s]", nic.name)
		utils.Assertf(nic.gateway != "", "gateway cannot be empty[nicname:%s]", nic.name)
		utils.Assertf(nic.netmask != "", "netmask cannot be empty[nicname:%s]", nic.name)
		utils.Assertf(nic.mac != "", "mac cannot be empty[nicname:%s]", nic.name)

		nicname, err := utils.GetNicNameByMac(nic.mac)
		utils.PanicOnError(err)
		if nicname != nic.name {
			devNames = append(devNames, &deviceName{
				expected: nic.name,
				actual:   nicname,
			})
		}
	}

	if len(devNames) != 0 {
		// shutdown links and change to temporary names
		cmds := make([]string, 0)
		for i, devname := range devNames {
			devnum := 1000 + i

			devname.swap = fmt.Sprintf("eth%v", devnum)
			cmds = append(cmds, fmt.Sprintf("ip link set dev %v down", devname.actual))
			cmds = append(cmds, fmt.Sprintf("ip link set dev %v name %v", devname.actual, devname.swap))
		}

		b := utils.Bash{
			Command: strings.Join(cmds, "\n"),
		}

		b.Run()
		b.PanicIfError()

		// change temporary names to real names and bring up links
		cmds = make([]string, 0)
		for _, devname := range devNames {
			cmds = append(cmds, fmt.Sprintf("ip link set dev %v name %v", devname.swap, devname.expected))
			cmds = append(cmds, fmt.Sprintf("ip link set dev %v up", devname.expected))
		}

		b = utils.Bash{
			Command: strings.Join(cmds, "\n"),
		}

		b.Run()
		b.PanicIfError()
	}

	vyos := vyos.NewParserFromShowConfiguration()
	tree := vyos.Tree

	/*
		sshkey := bootstrapInfo["publicKey"].(string)
		utils.Assert(sshkey != "", "cannot find 'publicKey' in bootstrap info")
		sshkeyparts := strings.Split(sshkey, " ")
		sshtype := sshkeyparts[0]
		key := sshkeyparts[1]
		id := sshkeyparts[2]

		tree.Setf("system login user vyos authentication public-keys %s key %s", id, key)
		tree.Setf("system login user vyos authentication public-keys %s type %s", id, sshtype)
	*/

	setNic := func(nic *nic) {
		cidr := utils.NetmaskToCIDR(nic.netmask)
		if cidr == -1 {
			panic(errors.New("netmask to cidr failed."))
		}

		//tree.Setf("interfaces ethernet %s hw-id %s", nic.name, nic.mac)
		tree.Setf("interfaces ethernet %s address %s", nic.name, fmt.Sprintf("%v/%v", nic.ip, cidr))
		tree.Setf("interfaces ethernet %s duplex auto", nic.name)
		tree.Setf("interfaces ethernet %s smp_affinity auto", nic.name)
		tree.Setf("interfaces ethernet %s speed auto", nic.name)
		if nic.isDefaultRoute {
			tree.Setf("system gateway-address %v", nic.gateway)
		}
	}

	/*
		sshport := bootstrapInfo["sshPort"].(float64)
		utils.Assert(sshport != 0, "sshport not found in bootstrap info")
	*/

	sshport := 22
	tree.Setf("service ssh port %v", int(sshport))
	tree.Setf("service ssh listen-address %v", eth0.ip)

	// configure firewall
	for _, nic := range nics {
		setNic(nic)

		tree.SetFirewallOnInterface(nic.name, "local",
			"action accept",
			"state established enable",
			"state related enable",
			fmt.Sprintf("destination address %v", nic.ip),
		)
		tree.SetFirewallOnInterface(nic.name, "local",
			"action accept",
			"protocol icmp",
			fmt.Sprintf("destination address %v", nic.ip),
		)

		tree.SetFirewallOnInterface(nic.name, "in",
			"action accept",
			"state established enable",
			"state related enable",
		)

		tree.SetFirewallWithRuleNumber(nic.name, "in", ROUTE_STATE_NEW_ENABLE_FIREWALL_RULE_NUMBER,
			"action accept",
			"state new enable",
		)

		tree.SetFirewallOnInterface(nic.name, "in",
			"action accept",
			"protocol icmp",
		)

		// only allow ssh traffic on eth0, disable on others
		if nic.name == "eth0" {
			tree.SetFirewallOnInterface(nic.name, "local",
				fmt.Sprintf("destination port %v", int(sshport)),
				fmt.Sprintf("destination address %v", nic.ip),
				"protocol tcp",
				"action accept",
			)
		} else {
			tree.SetFirewallOnInterface(nic.name, "local",
				fmt.Sprintf("destination port %v", int(sshport)),
				fmt.Sprintf("destination address %v", nic.ip),
				"protocol tcp",
				"action reject",
			)
		}

		tree.SetFirewallDefaultAction(nic.name, "local", "reject")
		tree.SetFirewallDefaultAction(nic.name, "in", "reject")

		tree.AttachFirewallToInterface(nic.name, "local")
		tree.AttachFirewallToInterface(nic.name, "in")
	}

	tree.Set("system time-zone Asia/Shanghai")

	/*
		password, found := bootstrapInfo["vyosPassword"]
		utils.Assert(found && password != "", "vyosPassword cannot be empty")
		tree.Setf("system login user vyos authentication plaintext-password %v", password)
	*/

	tree.Apply(true)

	arping := func(nicname, ip, gateway string) {
		b := utils.Bash{Command: fmt.Sprintf("arping -A -U -c 1 -I %s -s %s %s", nicname, ip, gateway)}
		b.Run()
	}

	// arping to advocate our mac addresses
	arping("eth0", eth0.ip, eth0.gateway)
	for _, nic := range nics {
		arping(nic.name, nic.ip, nic.gateway)
	}
}

func startZvr() {
	b := utils.Bash{
		Command: "bash -x /etc/init.d/zstack-virtualrouteragent restart >> /tmp/agentRestart.log 2>&1",
	}
	b.Run()
	b.PanicIfError()
}

func initDebugConfig() {
	octlog.InitDebugConfig(5)
}

func initLogConfig() {

	vyos.InitLog(5)
}

func initDebugAndLog() {
	initDebugConfig()
	initLogConfig()
}

func main() {

	initDebugAndLog()

	waitIptablesServiceOnline()

	waitVirtioPortOnline()
	parseKvmBootInfo()

	configureVyos()
	//	startZvr()
	octlog.Debug("successfully configured the sysmtem and bootstrap the octopuslink virtual router agents")
}