/***
Copyright 2014 Cisco Systems Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package comp

import (
	"github.com/contiv/netplugin/aci/model"
	"encoding/xml"
)

type Prov struct {
	model.BaseMo
	XMLName xml.Name `xml:"compProv"`
	Name string `xml:"name,attr,omitempty"`
	Vendor string  `xml:"pwd,domName,omitempty"`
	Ctrlrs []*Ctrlr `xml:"compCtrlr"`
	Doms []*Dom `xml:"compDom"`
}

func ProvDn(name string) model.Dn {
	return model.Dn("comp/prov-" + name)
}

func NewProv(name string) *Prov {
	return &Prov{
		BaseMo: model.BaseMo{Dn: ProvDn(name)},
		Name: name,
	}
}

type Dom struct {
	model.BaseMo
	XMLName xml.Name `xml:"compDom"`
	Name string `xml:"name,attr,omitempty"`
}

func DomDn(provName, domName string) model.Dn {
	return model.Dn(string(ProvDn(provName)) + "/dom-" + domName)
}

func NewDom(provName, domName string) *Dom {
	return &Dom{
		BaseMo: model.BaseMo{Dn: DomDn(provName, domName)},
		Name: domName,
	}
}

type Ctrlr struct {
	model.BaseMo
	XMLName xml.Name `xml:"compCtrlr"`
	Name string `xml:"name,attr,omitempty"`
	DomName string  `xml:"domName,attr,omitempty"`
	Hvs []*Hv `xml:"compHv"`
	Vms []*Vm `xml:"compVm"`
}

func CtrlrDn(provName, domName, ctrlrName string) model.Dn {
	return model.Dn(string(ProvDn(provName)) + "/ctrlr-[" + domName + "]-" + ctrlrName)
}

func NewCtrlr(provName, domName, ctrlrName string) *Ctrlr {
	return &Ctrlr{
		BaseMo: model.BaseMo{Dn: CtrlrDn(provName, domName, ctrlrName)},
		Name: ctrlrName,
		DomName: domName,
	}
}

func (c *Ctrlr) AddHv(oid string) *Hv {
	hv := &Hv{
		BaseMo: model.BaseMo{Dn:  model.Dn(string(c.Dn) + "/hv-" + oid)},
		Oid: oid,
		Name: oid,
	}
	c.Hvs = append(c.Hvs, hv)
	return hv
}

func (c *Ctrlr) AddVm(hv *Hv, oid string) *Vm {
	vm := &Vm{
		BaseMo: model.BaseMo{Dn:  model.Dn(string(c.Dn) + "/vm-" + oid)},
		RsHv: &RsHv{TDn: hv.Dn},
		Oid: oid,
		Name: oid,
	}
	c.Vms = append(c.Vms, vm)
	return vm
}

type Hv struct {
	model.BaseMo
	XMLName xml.Name `xml:"compHv"`
	Name string `xml:"name,attr,omitempty"`
	Oid string `xml:"oid,attr,omitempty"`
	Type string `xml:"type,attr,omitempty"`
	HpNics []*HpNic `xml:"compHpNic"`
}

func (hv *Hv) AddHpNic(oid string) *HpNic {
	nic := &HpNic {
		BaseMo: model.BaseMo{Dn:  model.Dn(string(hv.Dn) + "/hpnic-" + oid)},
		Oid: oid,
		Name: oid,
	}
	hv.HpNics = append(hv.HpNics, nic)
	return nic
}

type HpNic struct {
	model.BaseMo
	XMLName xml.Name `xml:"compHpNic"`
	Name string `xml:"name,attr,omitempty"`
	Oid string `xml:"oid,attr,omitempty"`
}

type Vm struct {
	model.BaseMo
	XMLName xml.Name `xml:"compVm"`
	Name string `xml:"name,attr,omitempty"`
	Oid string `xml:"oid,attr,omitempty"`
	RsHv *RsHv `xml:"compRsHv"`
	VNics []*VNic `xml:"compVNic"`
}

func (vm *Vm) AddVNic(mac string) *VNic {
	nic := &VNic {
		BaseMo: model.BaseMo{Dn:  model.Dn(string(vm.Dn) + "/vnic-" + mac)},
		Mac: mac,
	}
	vm.VNics = append(vm.VNics, nic)
	return nic
}

type RsHv struct {
	model.BaseMo
	XMLName xml.Name `xml:"compRsHv"`
	Name string `xml:"name,attr,omitempty"`
	TDn model.Dn `xml:"tDn,attr,omitempty"`
}

type VNic struct {
	model.BaseMo
	XMLName xml.Name `xml:"compVNic"`
	Name string `xml:"name,attr,omitempty"`
	Oid string `xml:"oid,attr,omitempty"`
	Mac string `xml:"mac,attr,omitempty"`
	Type string `xml:"type,attr,omitempty"`
	OperSt string `xml:"operSt,attr,omitempty"`
	Ips []*Ip `xml:"compIp"`
}


func (nic *VNic) AddIp(addr string) *Ip {
	ip := &Ip {
		BaseMo: model.BaseMo{Dn:  model.Dn(string(nic.Dn) + "/ip-" + addr)},
		Addr: addr,
	}
	nic.Ips = append(nic.Ips, ip)
	return ip
}

type Ip struct {
	model.BaseMo
	XMLName xml.Name `xml:"compIp"`
	Addr string `xml:"addr,attr,omitempty"`
}


