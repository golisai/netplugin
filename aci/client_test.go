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

package aci

import (
	"encoding/xml"
	"flag"
	"github.com/contiv/netplugin/aci/model/aaa"
	"github.com/contiv/netplugin/aci/model/comp"
	"github.com/contiv/netplugin/aci/model/fv"
	"log"
	"os"
	"testing"
)

var client *Client

func TestDnQuery(t *testing.T) {
	var user aaa.User
	if err := client.DnQuery("uni/userext/user-" + client.User, &user, false); err != nil {
		t.Fatal(err)
	}

	if user.Name != client.User {
		data, _ := xml.Marshal(user)
		t.Log(string(data))
		t.Error("aaaUser DN query failed")
	}
}

func TestClassQuery(t *testing.T) {
	var tenants []fv.Tenant
	if err := client.ClassQuery(&tenants, false); err != nil {
		t.Fatal(err)
	}
	
	// minimum common, infra, mgmt tenants should be there 
	if len(tenants) < 3 {
		for _,tn := range tenants {
			data, _ := xml.Marshal(tn)
			t.Log(string(data))
		}
		
		t.Errorf("Tenant class query returned only %v mos", len(tenants))
	}
}

func TestVmmModel(t *testing.T) {
		model := `
<compProv name="Microsoft">
  <compCtrlr name="contiv" domName="test">
    <compHv name="hv" oid="hv" type="hv">
      <compHpNic name="pnic1" oid="pnic1"/>
    </compHv>
    <compVm name="vm1" oid="vm1">
       <compRsHv tDn="comp/prov-VMWare/ctrlr-[test]-contiv/hv-hv"/>
       <compVNic name="vinc1" oid="vnic1" mac="AA:AA:AA:AA:AA:AA" operSt="up">
          <compIp addr="1.1.1.1"/>
       </compVNic>
    </compVm>
  </compCtrlr>
</compProv>
`
	var prov comp.Prov
	xml.Unmarshal([]byte(model), &prov)

	if err := client.Post(client.URL + "/testapi/mo/comp/prov-Microsoft.xml", &prov); err != nil {
		t.Fatal(err)
	}
	
	var mos []comp.Ctrlr
	if err := client.ClassQuery(&mos, true); err != nil {
		t.Fatal(err)
	}

	data,_ := xml.Marshal(mos)
	t.Log(string(data))
}

func TestMain(m *testing.M) {
	url := flag.String("aci.url", "http://srgoli-bld:8000", "APIC Rest URL")
	user := flag.String("aci.user", "admin", "APIC login username")
	passwd := flag.String("aci.passwd", "ins3965!", "APIC login passwd")
	
	flag.Parse()

	var err error
	client, err = NewClient(*url, *user, *passwd)
	if err != nil {
		log.Fatal("Cant' create client", err)
	}
	os.Exit(m.Run())
}
