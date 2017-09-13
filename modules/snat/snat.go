package snat

import (
	"fmt"
	"octlink/ovs/utils"
	"octlink/ovs/utils/merrors"
	"octlink/ovs/utils/octlog"
	"octlink/ovs/utils/vyos"
)

var logger *octlog.LogConfig

// InitLog to init log config
func InitLog(level int) {
	logger = octlog.InitLogConfig("image.log", level)
}

const (
	// SetSnatPath for Set SNAT
	SetSnatPath = "/setsnat"

	// RemoveSnatPath for remove snat path
	RemoveSnatPath = "/removesnat"

	// SyncSnatPath for sync snat
	SyncSnatPath = "/syncsnat"
)

// SnatRuleNumber for max snat rule number
var SnatRuleNumber = 9999

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

// Add for image, after image added,
// installpath, diskSize, virtualSize, Status, md5sum need update after manifest installed
func (s *Snat) Add() int {

	tree := vyos.NewParserFromShowConfiguration().Tree
	outNic, err := utils.GetNicNameByMac(s.PublicNicMac)

	utils.PanicOnError(err)

	address, err := utils.GetNetworkNumber(s.PrivateNicIP, s.SnatNetmask)
	utils.PanicOnError(err)

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

// Remove Snat rule
func (s *Snat) Remove() int {

	tree := vyos.NewParserFromShowConfiguration().Tree
	rs := tree.Get("nat source rule")
	if rs == nil {
		logger.Debugf("not nat source rule remove\n")
		return merrors.ErrSuccess
	}

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

// Sync Snat rule, delete it firstly and then add it back.
func (s *Snat) Sync() int {

	tree := vyos.NewParserFromShowConfiguration().Tree

	outNic, err := utils.GetNicNameByMac(s.PublicNicMac)
	if err != nil {
		logger.Panicf("Get Nic Name by Mac %s error %s\n", s.PublicNicMac, err)
		return merrors.ErrSegmentNotExist
	}

	address, err := utils.GetNetworkNumber(s.PrivateNicIP, s.SnatNetmask)
	if err != nil {
		logger.Panicf("Get Network Number error %s:%s %s\n",
			s.PrivateNicIP, s.SnatNetmask, err)
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
func GetSnat(privateIP string, netmask string) *Snat {
	return nil
}

// GetAllSnats by condition
func GetAllSnats() []*Snat {

	tree := vyos.NewParserFromShowConfiguration().Tree

	if rs := tree.Getf("nat source rule"); rs != nil {
		logger.Debugf("got nat rule of %s\n", rs.String())
	}

	return make([]*Snat, 0)
}
