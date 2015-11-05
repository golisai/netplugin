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
	model.Mo
	XMLName xml.Name `xml:"compProv"`
	Name string `xml:"name,attr,omitempty"`
	Vendor string  `xml:"pwd,domName,omitempty"`
	Ctrlrs []*Ctrlr `xml:"compCtrlr"`
	Doms []*Dom `xml:"compDom"`
}

type Dom struct {
	model.Mo
	XMLName xml.Name `xml:"compDom"`
	Name string `xml:"name,attr,omitempty"`
}

type Ctrlr struct {
	model.Mo
	XMLName xml.Name `xml:"compCtrlr"`
	Name string `xml:"name,attr,omitempty"`
	DomName string  `xml:"domName,attr,omitempty"`
	Hvs []*Hv `xml:"compHv"`
	Vms []*Vm `xml:"compVm"`
}

type Hv struct {
	model.Mo
	XMLName xml.Name `xml:"compHv"`
	Name string `xml:"name,attr,omitempty"`
	Oid string `xml:"oid,attr,omitempty"`
	Type string `xml:"type,attr,omitempty"`
	HpNics []*HpNic `xml:"compHpNic"`
}

type HpNic struct {
	model.Mo
	XMLName xml.Name `xml:"compHpNic"`
	Name string `xml:"name,attr,omitempty"`
	Oid string `xml:"oid,attr,omitempty"`
}

type Vm struct {
	model.Mo
	XMLName xml.Name `xml:"compVm"`
	Name string `xml:"name,attr,omitempty"`
	Oid string `xml:"oid,attr,omitempty"`
	RsHv *RsHv `xml:"compRsHv"`
	VNics []*VNic `xml:"compVNic"`
}

type RsHv struct {
	model.Mo
	XMLName xml.Name `xml:"compRsHv"`
	Name string `xml:"name,attr,omitempty"`
	TDn string `xml:"tDn,attr,omitempty"`
}

type VNic struct {
	model.Mo
	XMLName xml.Name `xml:"compVNic"`
	Name string `xml:"name,attr,omitempty"`
	Oid string `xml:"oid,attr,omitempty"`
	Mac string `xml:"mac,attr,omitempty"`
	Type string `xml:"type,attr,omitempty"`
	OperSt string `xml:"operSt,attr,omitempty"`
	Ips []*Ip `xml:"compIp"`
}

type Ip struct {
	model.Mo
	XMLName xml.Name `xml:"compIp"`
	Addr string `xml:"addr,attr,omitempty"`
}


