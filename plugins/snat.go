package plugins

import (
	"fmt"
	"octlink/ovs/utils"
	"octlink/ovs/utils/merrors"
	"octlink/ovs/utils/vyos"
)

// Snat for snat sturcture
type Snat struct {
	PrivateNicMac string `json:"privateNicMac"`
	PrivateNicIP  string `json:"privateNicIp"`
	PublicIP      string `json:"publicIp"`
	PublicNicMac  string `json:"publicNicMac"`
	SnatNetmask   string `json:"snatNetmask"`
}

// GetSnatCount to return image count by condition
func GetSnatCount() int {
	return len(GetAllSnats())
}

func hasRuleNumberForAddress(tree *vyos.ConfigTree, address string) bool {

	rs := tree.Get("nat source rule")
	if rs == nil {
		return false
	}

	for _, r := range rs.Children() {
		if addr := r.Get("source address"); addr != nil && addr.Value() == address {
			return true
		}
	}

	return false
}

// AddSnat for image, after image added,
// installpath, diskSize, virtualSize, Status, md5sum need update after manifest installed
func (s *Snat) AddSnat() int {

	tree := vyos.NewParserFromShowConfiguration().Tree

	outNic, err := utils.GetNicNameByMac(s.PublicNicMac)
	if err != nil {
		logger.Panicf("get nic name by mac %s error %s\n", s.PublicNicMac, err)
		return merrors.ErrBadParas
	}

	privateIP, snatNetmask, _, err := utils.GetNicInfoByMac(s.PrivateNicMac)
	if err != nil {
		logger.Panicf("get nic info by mac %s error, %s\n", s.PrivateNicMac, err)
		return merrors.ErrBadParas
	}

	address, err := utils.GetNetworkNumber(privateIP, snatNetmask)
	if err != nil {
		logger.Panicf("get network number by %s:%s error, %s\n",
			privateIP, snatNetmask, err)
		return merrors.ErrBadParas
	}

	if hasRuleNumberForAddress(tree, address) {
		logger.Errorf("not enough rule number for snat, address[%s]\n", address)
		return merrors.ErrSyscallErr
	}

	// make source nat rule as the latest rule
	// in case there are EIP rules
	tree.SetSnatWithRuleNumber(SnatRuleNumber,
		fmt.Sprintf("outbound-interface %s", outNic),
		fmt.Sprintf("source address %v", address),
		fmt.Sprintf("translation address %s", s.PublicIP),
	)

	tree.Apply(false)

	return 0
}

// RemoveSnat Snat rule
func (s *Snat) RemoveSnat() int {

	tree := vyos.NewParserFromShowConfiguration().Tree
	rs := tree.Get("nat source rule")
	if rs == nil {
		logger.Debugf("not nat source rule remove\n")
		return merrors.ErrSuccess
	}

	privateNicIP, snatNetmask, _, err := utils.GetNicInfoByMac(s.PrivateNicMac)
	if err != nil {
		logger.Errorf("get nic info of %s error\n", s.PrivateNicMac)
		return merrors.ErrBadParas
	}

	s.PrivateNicIP = privateNicIP
	s.SnatNetmask = snatNetmask

	logger.Debugf("get network of %s source rule of %s:%s\n",
		s.PrivateNicMac, s.PrivateNicIP, s.SnatNetmask)

	address, err := utils.GetNetworkNumber(s.PrivateNicIP, s.SnatNetmask)
	if err != nil {
		logger.Panicf("nat source rule of %s:%s not exist\n", s.PrivateNicIP, s.SnatNetmask)
		return merrors.ErrBadParas
	}

	for _, r := range rs.Children() {
		if addr := r.Get("source address"); addr != nil && addr.Value() == address {
			addr.Delete()
		}
	}

	tree.Apply(false)

	return merrors.ErrSuccess
}

// SyncSnat Snat rule, delete it firstly and then add it back.
func (s *Snat) SyncSnat() int {

	tree := vyos.NewParserFromShowConfiguration().Tree

	outNic, err := utils.GetNicNameByMac(s.PublicNicMac)
	if err != nil {
		logger.Panicf("get nic name by mac %s error %s\n", s.PublicNicMac, err)
		return merrors.ErrBadParas
	}

	privateIP, snatNetmask, _, err := utils.GetNicInfoByMac(s.PrivateNicMac)
	if err != nil {
		logger.Panicf("get nic info by mac %s error, %s\n", s.PrivateNicMac, err)
		return merrors.ErrBadParas
	}

	address, err := utils.GetNetworkNumber(privateIP, snatNetmask)
	if err != nil {
		logger.Panicf("get network number by %s:%s error, %s\n",
			privateIP, snatNetmask, err)
		return merrors.ErrBadParas
	}

	if rs := tree.Getf("nat source rule %v", SnatRuleNumber); rs != nil {
		rs.Delete()
	}

	tree.SetSnatWithRuleNumber(SnatRuleNumber,
		fmt.Sprintf("outbound-interface %s", outNic),
		fmt.Sprintf("source address %s", address),
		fmt.Sprintf("translation address %s", s.PublicIP),
	)

	tree.Apply(false)

	return 0
}

// GetSnat get snat settings
func GetSnat(privateNicMac string) (*Snat, int) {

	privateIP, netmask, _, err := utils.GetNicInfoByMac(privateNicMac)
	if err != nil {
		logger.Errorf("not nic info got for private mac %s\n", privateNicMac)
		return nil, merrors.ErrBadParas
	}
	network, _ := utils.GetNetworkNumber(privateIP, netmask)

	rules := GetAllSnats()
	for _, nat := range rules {
		if n, _ := utils.GetNetworkNumber(nat.PrivateNicIP, nat.SnatNetmask); n == network {
			logger.Debugf("found nat rule of %s\n", network)
			return nat, merrors.ErrSuccess
		}
	}

	logger.Errorf("not found nat rule of %s\n", network)

	return nil, merrors.ErrSegmentNotExist
}

// GetAllSnats by condition
func GetAllSnats() []*Snat {

	tree := vyos.NewParserFromShowConfiguration().Tree

	sn := new(Snat)

	rule := tree.Getf("nat source rule %d", SnatRuleNumber)
	if rule != nil {
		outNic := rule.Get("outbound-interface").Value()
		sn.PublicNicMac = utils.GetNicMacByName(outNic)
		logger.Debugf("Got nat oubound-interface %s:%s\n", outNic, sn.PublicNicMac)
	}

	if rs := rule.Getf("source address"); rs != nil {
		addr, netmask := utils.ParseCIDR(rs.Value())
		sn.PrivateNicIP = addr
		sn.SnatNetmask = netmask
		logger.Debugf("Got nat private nic ip %s/%s\n", sn.PrivateNicIP, sn.SnatNetmask)
	}

	if rs := rule.Getf("translation address"); rs != nil {
		sn.PublicIP = rs.Value()
		logger.Debugf("Got nat public nic ip %s\n", sn.PrivateNicIP)
	}

	/*

		if rs := tree.Getf("nat source rule %d source", SnatRuleNumber); rs != nil {
			addr, netmask := utils.ParseCIDR(rs.Get("address").Value())
			sn.PrivateNicIP = addr
			sn.SnatNetmask = netmask
			logger.Debugf("Got nat private nic ip %s/%s\n", sn.PrivateNicIP, sn.SnatNetmask)
		}

		if rs := tree.Getf("nat source rule %d translation", SnatRuleNumber); rs != nil {
			sn.PublicIP = rs.Get("address").Value()
			logger.Debugf("Got nat public nic ip %s\n", sn.PrivateNicIP)
		}
	*/

	return []*Snat{
		sn,
	}
}
