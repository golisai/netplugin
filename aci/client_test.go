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
	"github.com/contiv/netplugin/aci/model/fv"
	"log"
	"os"
	"testing"
)

var client *Client

func TestDnQuery(t *testing.T) {
	var user aaa.User
	if err := client.DnQuery("uni/userext/user-" + client.User, &user); err != nil {
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
	if err := client.ClassQuery(&tenants); err != nil {
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
