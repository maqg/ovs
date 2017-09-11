package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
)

const (
	RVM_ROUTE_PROTO            = "rvm"
	RVM_ROUTE_PROTO_IDENTIFFER = "192"
)

func NetmaskToCIDR(netmask string) (int, error) {
	countBit := func(num uint) int {
		count := uint(0)
		var i uint
		for i = 31; i > 0; i-- {
			count += ((num << i) >> uint(31)) & uint(1)
		}

		return int(count)
	}

	cidr := 0
	for _, o := range strings.Split(netmask, ".") {
		num, err := strconv.ParseUint(o, 10, 32)
		if err != nil {
			return -1, err
		}
		cidr += countBit(uint(num))
	}

	return cidr, nil
}

func GetNetworkNumber(ip, netmask string) (string, error) {
	ips := strings.Split(ip, ".")
	masks := strings.Split(netmask, ".")

	ipInByte := make([]interface{}, 4)
	for i := 0; i < len(ips); i++ {
		p, err := strconv.ParseUint(ips[i], 10, 32)
		if err != nil {
			return "", errors.Wrap(err, fmt.Sprintf("unable to get network number[ip:%v, netmask:%v]", ip, netmask))
		}
		m, err := strconv.ParseUint(masks[i], 10, 32)
		PanicOnError(err)
		if err != nil {
			return "", errors.Wrap(err, fmt.Sprintf("unable to get network number[ip:%v, netmask:%v]", ip, netmask))
		}
		ipInByte[i] = p & m
	}

	cidr, err := NetmaskToCIDR(netmask)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("unable to get network number[ip:%v, netmask:%v]", ip, netmask))
	}

	return fmt.Sprintf("%v.%v.%v.%v/%v", ipInByte[0], ipInByte[1], ipInByte[2], ipInByte[3], cidr), nil
}

type Nic struct {
	Name string
	Mac  string
}

func (nic Nic) String() string {
	s, _ := json.Marshal(nic)
	return string(s)
}

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

func GetNicNameByIp(ip string) (string, error) {
	bash := Bash{
		Command: fmt.Sprintf("ip addr | grep -w %s", ip),
	}
	ret, o, _, err := bash.RunWithReturn()
	if err != nil {
		return "", err
	}
	if ret != 0 {
		return "", errors.New(fmt.Sprintf("no nic with the IP[%s] found in the system", ip))
	}

	o = strings.TrimSpace(o)
	os := strings.Split(o, " ")
	return os[len(os)-1], nil
}

func GetIpFromUrl(url string) (string, error) {
	ip := strings.Split(strings.Split(url, "/")[2], ":")[0]
	return ip, nil
}

func SetZStackRoute(ip string, nic string) error {
	SetZStackRouteProtoIdentifier()
	bash := Bash{
		Command: fmt.Sprintf("ip route add %s/32 dev %s proto %s", ip, nic, RVM_ROUTE_PROTO),
	}
	ret, _, _, err := bash.RunWithReturn()
	if err != nil {
		return err
	}
	// NOTE(WeiW): It will return 2 if exists
	if ret != 0 && ret != 2 {
		return errors.New(fmt.Sprintf("add route to %s/32 use dev %s failed", ip, nic))
	}

	return nil
}

func RemoveZStackRoute(ip string, nic string) error {
	SetZStackRouteProtoIdentifier()
	bash := Bash{
		Command: fmt.Sprintf("ip route del %s/32 dev %s proto %s", ip, nic, RVM_ROUTE_PROTO),
	}
	ret, _, _, err := bash.RunWithReturn()
	if err != nil {
		return err
	}
	if ret != 0 {
		return errors.New(fmt.Sprintf("del route to %s/32 use dev %s failed", ip, nic))
	}

	return nil
}

func SetZStackRouteProtoIdentifier() {
	bash := Bash{
		Command: "grep rvm /etc/iproute2/rt_protos",
	}
	check, _, _, _ := bash.RunWithReturn()

	if check != 0 {
		log.Debugf("no route proto rvm in /etc/iproute2/rt_protos")
		bash = Bash{
			Command: fmt.Sprintf("sudo bash -c \"echo -e '\n\n# Used by rvm\n%s     rvm' >> /etc/iproute2/rt_protos\"", RVM_ROUTE_PROTO_IDENTIFFER),
		}
		bash.Run()
	}
}
