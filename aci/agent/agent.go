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

package main

import (
	"flag"
	"github.com/contiv/netplugin/aci"
	"github.com/contiv/netplugin/aci/model/comp"
	"github.com/samalba/dockerclient"
	"log"
	"os"
	"strings"
)

const Provider = "VMware"

func addContainer(clt *aci.Client, info *dockerclient.ContainerInfo) error {
	log.Println("Adding container: ", info.Name, info.NetworkSettings.IPAddress)
	// Create Ctrlr
	ctrlr := comp.NewCtrlr(Provider, "contiv", "contiv")
	
	// Add HV
	hostName, _ := os.Hostname()
	hv := ctrlr.AddHv(hostName)
	hv.AddHpNic("eth0")

	name := strings.TrimPrefix(info.Name, "/")
	
	vm := ctrlr.AddVm(hv, name)
	vnic := vm.AddVNic(info.Config.MacAddress)
	vnic.AddIp(info.NetworkSettings.IPAddress)

	return clt.Write(ctrlr)
}

func main() {
	url := flag.String("url", "http://srgoli-bld:8000", "APIC Rest URL")
	user := flag.String("user", "admin", "APIC login username")
	passwd := flag.String("passwd", "ins3965!", "APIC login passwd")

	flag.Parse()

	clt, err := aci.NewClient(*url, *user, *passwd)
	if err != nil {
		log.Fatal("Can't create client", err)
	}

	// Init the client
	docker, err := dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)
	if err != nil {
		log.Fatal("Cant' create docker client", err)
	}
	
	// Get running containers
	containers, err := docker.ListContainers(false, false, "")
	if err != nil {
		log.Fatal("Can't list containers", err)
	}
	
	for _, c := range containers {
		info, err := docker.InspectContainer(c.Id)
		if err != nil {
			log.Fatal("Inspect failed", err)	
		}
		if err = addContainer(clt, info); err != nil {
			log.Fatal("addContainer failed", err)
		}
	}
	
}
