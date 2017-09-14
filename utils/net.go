package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
)

const (
	// RvmRouteProto rvm route proto
	RvmRouteProto = "rvm"

	// RvmRouteProtoIdentifier with 192
	RvmRouteProtoIdentifier = "192"
)

// Nic base structure
type Nic struct {
	Name string
	Mac  string
}

// String convert Nic to string
func (nic Nic) String() string {
	s, _ := json.Marshal(nic)
	return string(s)
}

// GetAllNics with name and mac address
func GetAllNics() (map[string]Nic, error) {
	const ROOT = "/sys/class/net"

	files, err := ioutil.ReadDir(ROOT)
	if err != nil {
		return nil, err
	}

	nics := make(map[string]Nic)
	for _, f := range files {
		if f.IsDir() || f.Name() == "lo" {
			continue
		}

		macfile := filepath.Join(ROOT, f.Name(), "address")
		mac, err := ioutil.ReadFile(macfile)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("unable to read the mac file[%s]", macfile))
		}
		nics[f.Name()] = Nic{
			Name: strings.TrimSpace(f.Name()),
			Mac:  strings.TrimSpace(string(mac)),
		}
	}

	return nics, nil
}

// GetNicNameByMac get nicname by mac address
func GetNicNameByMac(mac string) (string, error) {
	nics, err := GetAllNics()
	if err != nil {
		return "", err
	}

	for _, nic := range nics {
		if nic.Mac == mac {
			return nic.Name, nil
		}
	}

	return "", fmt.Errorf("cannot find any nic with the mac[%s]", mac)
}

// GetNicMacByName get mac address by nic name
func GetNicMacByName(nicname string) string {
	nics, _ := GetAllNics()
	for _, nic := range nics {
		if nic.Name == nicname {
			return nic.Mac
		}
	}

	return ""
}

// GetNicIP get nic IP by nicname
func GetNicIP(nicname string) string {
	ip, _, _, _ := GetNicInfo(nicname)
	return ip
}

// GetNicIPByMac get nic IP by mac
func GetNicIPByMac(mac string) string {

	nicname, err := GetNicNameByMac(mac)
	if err != nil {
		return ""
	}

	ip, _, _, _ := GetNicInfo(nicname)

	return ip
}

// GetNicInfoByMac get nic info by mac address
func GetNicInfoByMac(mac string) (string, string, string, error) {
	nicname, err := GetNicNameByMac(mac)
	if err != nil {
		return "", "", "", err
	}

	return GetNicInfo(nicname)
}

// GetNicInfo get ip address,netmask and network by nic name
func GetNicInfo(nicname string) (string, string, string, error) {
	bash := Bash{
		Command: fmt.Sprintf("ip addr show %s | grep -w inet", nicname),
	}
	ret, o, _, err := bash.RunWithReturn()
	if err != nil {
		return "", "", "", err
	}
	if ret != 0 {
		return "", "", "", fmt.Errorf("no nic info the name of [%s] found in the system", nicname)
	}

	o = strings.TrimSpace(o)
	os := strings.Split(o, " ")

	addr := strings.Split(os[1], "/")

	return addr[0], CIDRToNetmask(StringToInt(addr[1])), os[3], nil
}

// GetNicNameByIP get nic name by ip address
func GetNicNameByIP(ip string) (string, error) {
	bash := Bash{
		Command: fmt.Sprintf("ip addr | grep -w %s", ip),
	}
	ret, o, _, err := bash.RunWithReturn()
	if err != nil {
		return "", err
	}
	if ret != 0 {
		return "", fmt.Errorf("no nic with the IP[%s] found in the system", ip)
	}

	o = strings.TrimSpace(o)
	os := strings.Split(o, " ")
	return os[len(os)-1], nil
}

// GetIPFromURL get ip address from url
func GetIPFromURL(url string) (string, error) {
	ip := strings.Split(strings.Split(url, "/")[2], ":")[0]
	return ip, nil
}

// SetRvmRoute to set rvm route config
func SetRvmRoute(ip string, nic string) error {
	SetRvmRouteProtoIdentifier()
	bash := Bash{
		Command: fmt.Sprintf("ip route add %s/32 dev %s proto %s", ip, nic, RvmRouteProto),
	}
	ret, _, _, err := bash.RunWithReturn()
	if err != nil {
		return err
	}
	// NOTE(WeiW): It will return 2 if exists
	if ret != 0 && ret != 2 {
		return fmt.Errorf("add route to %s/32 use dev %s failed", ip, nic)
	}

	return nil
}

// RemoveRvmRoute to remove rvm route config
func RemoveRvmRoute(ip string, nic string) error {
	SetRvmRouteProtoIdentifier()
	bash := Bash{
		Command: fmt.Sprintf("ip route del %s/32 dev %s proto %s", ip, nic, RvmRouteProto),
	}
	ret, _, _, err := bash.RunWithReturn()
	if err != nil {
		return err
	}
	if ret != 0 {
		return fmt.Errorf("del route to %s/32 use dev %s failed", ip, nic)
	}

	return nil
}

// SetRvmRouteProtoIdentifier for rvm
func SetRvmRouteProtoIdentifier() {
	bash := Bash{
		Command: "grep rvm /etc/iproute2/rt_protos",
	}
	check, _, _, _ := bash.RunWithReturn()

	if check != 0 {
		log.Debugf("no route proto rvm in /etc/iproute2/rt_protos")
		bash = Bash{
			Command: fmt.Sprintf("sudo bash -c \"echo -e '\n\n# Used by rvm\n%s     rvm' >> /etc/iproute2/rt_protos\"", RvmRouteProtoIdentifier),
		}
		bash.Run()
	}
}
