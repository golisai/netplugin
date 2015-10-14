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

package singlehost

import (
	"strings"
	"testing"
	"time"

	"github.com/contiv/netplugin/systemtests/utils"
	u "github.com/contiv/netplugin/utils"
)

// Testcase:
// - Create a single vlan network with two endpoints
// - Verify that the endpoints are able to ping
func TestSingleHostSingleVlanPingSuccess_sanity(t *testing.T) {
	defer func() {
		utils.ConfigCleanupCommon(t, testbed.GetNodes())
		utils.StopOnError(t.Failed())
	}()

	//create a single vlan network, with two endpoints
	jsonCfg :=
		`{
        "Tenants" : [ {
            "Name"                      : "tenant-one",
            "DefaultNetType"            : "vlan",
            "SubnetPool"                : "11.1.0.0/16",
            "AllocSubnetLen"            : 24,
            "Vlans"                     : "11-48",
            "Networks"  : [ {
                "Name"                  : "orange",
                "Endpoints" : [
                {
                    "Container"         : "myContainer1",
                    "Host"              : "host1"
                },
                {
                    "Container"         : "myContainer2",
                    "Host"              : "host1"
                } ]
            } ]
        } ]
        }`

	utils.ConfigSetupCommon(t, jsonCfg, testbed.GetNodes())

	node := testbed.GetNodes()[0]

	utils.StartServer(t, node, "myContainer1")
	defer func() {
		utils.DockerCleanup(t, node, "myContainer1")
	}()

	ipAddress := utils.GetIPAddress(t, node, "orange-myContainer1", u.EtcdNameStr)
	utils.StartClient(t, node, "myContainer2", ipAddress)
	defer func() {
		utils.DockerCleanup(t, node, "myContainer2")
	}()
}

// Testcase:
// - Create a network with two vlans with two endpoints each
// - Verify that the endpoints in same vlan are able to ping
func TestSingleHostMultiVlanPingSuccess_sanity(t *testing.T) {
	defer func() {
		utils.ConfigCleanupCommon(t, testbed.GetNodes())
		utils.StopOnError(t.Failed())
	}()

	//create a single vlan network, with two endpoints
	jsonCfg :=
		`{
        "Tenants" : [ {
            "Name"                      : "tenant-one",
            "DefaultNetType"            : "vlan",
            "SubnetPool"                : "11.1.0.0/16",
            "AllocSubnetLen"            : 24,
            "Vlans"                     : "11-48",
            "Networks"  : [ {
                "Name"                  : "orange",
                "Endpoints" : [
                {
                    "Container"         : "myContainer1",
                    "Host"              : "host1"
                },
                {
                    "Container"         : "myContainer2",
                    "Host"              : "host1"
                } ]
            },
            {
                "Name"                  : "purple",
                "Endpoints" : [
                {
                    "Container"         : "myContainer3",
                    "Host"              : "host1"
                },
                {
                    "Container"         : "myContainer4",
                    "Host"              : "host1"
                } ]
            } ]
        } ]
        }`

	utils.ConfigSetupCommon(t, jsonCfg, testbed.GetNodes())

	node := testbed.GetNodes()[0]
	utils.StartServer(t, node, "myContainer1")
	defer func() {
		utils.DockerCleanup(t, node, "myContainer1")
	}()

	ipAddress := utils.GetIPAddress(t, node, "orange-myContainer1", u.EtcdNameStr)
	utils.StartClient(t, node, "myContainer2", ipAddress)
	defer func() {
		utils.DockerCleanup(t, node, "myContainer2")
	}()

	utils.StartServer(t, node, "myContainer4")
	defer func() {
		utils.DockerCleanup(t, node, "myContainer4")
	}()

	ipAddress = utils.GetIPAddress(t, node, "purple-myContainer4", u.EtcdNameStr)
	utils.StartClient(t, node, "myContainer3", ipAddress)
	defer func() {
		utils.DockerCleanup(t, node, "myContainer3")
	}()
}

// Testcase:
// - Create a network with two vlans with one endpoints each
// - Verify that the endpoints in different vlans are not able to ping
func TestSingleHostMultiVlanPingFailure_sanity(t *testing.T) {
	defer func() {
		utils.ConfigCleanupCommon(t, testbed.GetNodes())
		utils.StopOnError(t.Failed())
	}()

	//create a single vlan network, with two endpoints
	jsonCfg :=
		`{
        "Tenants" : [ {
            "Name"                      : "tenant-one",
            "DefaultNetType"            : "vlan",
            "SubnetPool"                : "11.1.0.0/16",
            "AllocSubnetLen"            : 24,
            "Vlans"                     : "11-48",
            "Networks"  : [ {
                "Name"                  : "orange",
                "Endpoints" : [
                {
                    "Container"         : "myContainer1",
                    "Host"              : "host1"
                } ]
            },
            {
                "Name"                  : "purple",
                "Endpoints" : [
                {
                    "Container"         : "myContainer2",
                    "Host"              : "host1"
                } ]
            } ]
        } ]
        }`

	utils.ConfigSetupCommon(t, jsonCfg, testbed.GetNodes())

	node := testbed.GetNodes()[0]

	utils.StartServer(t, node, "myContainer1")
	defer func() {
		utils.DockerCleanup(t, node, "myContainer1")
	}()

	ipAddress := utils.GetIPAddress(t, node, "orange-myContainer1", u.EtcdNameStr)
	utils.StartClientFailure(t, node, "myContainer2", ipAddress)
	defer func() {
		utils.DockerCleanup(t, node, "myContainer2")
	}()
}

// Default Network Assignment and freeing
func TestSingleHostDefaultNetwork_sanity(t *testing.T) {
	defer func() {
		utils.ConfigCleanupCommon(t, testbed.GetNodes())
		utils.StopOnError(t.Failed())
	}()

	//create a single vlan network, with two endpoints
	jsonCfg :=
		`{
        "Tenants" : [ {
            "Name"                      : "tee-one",
            "DefaultNetType"            : "vlan",
            "DefaultNetwork"            : "orange",
            "SubnetPool"                : "10.240.0.0/16",
            "AllocSubnetLen"            : 24,
            "Vlans"                     : "11-48",
            "Networks"  : [ {
                "Name"                  : "orange",
                "Endpoints" : [
                {
                    "Container"         : "myContainer1",
                    "Host"              : "host1"
                } ]
            },
            {
                "Name"                  : "purple",
                "Endpoints" : [
                {
                    "Container"         : "myContainer2",
                    "Host"              : "host1"
                } ]
            } ]
        } ]
        }`

	utils.ConfigSetupCommon(t, jsonCfg, testbed.GetNodes())

	node := testbed.GetNodes()[0]

	utils.StartServer(t, node, "myContainer1")
	defer func() {
		utils.DockerCleanup(t, node, "myContainer1")
	}()

	ipAddress := utils.GetIPAddress(t, node, "orange-myContainer1", u.EtcdNameStr)
	utils.StartClientFailure(t, node, "myContainer2", ipAddress)
	defer func() {
		utils.DockerCleanup(t, node, "myContainer2")
	}()

	// confirm the default gateway in one of the containers
	output, err := node.RunCommandWithOutput("docker exec myContainer1 /sbin/ip route")
	if err != nil {
		t.Fatalf("Error - unable to get default ip route, output = '%s'", output)
	}
	if !strings.Contains(output, "default via 10.240.0.254") {
		t.Fatalf("Error - unable to confirm container's default ip route, output = '%s'", output)
	}

	jsonCfg = `
    {
      "Tenants" : [ {
        "Name"                      : "tee-one",
        "Networks"  : [ {
          "Name"                    : "orange",
          "Endpoints" : [ {
            "Container"             : "myContainer1",
            "Host"                  : "host1"
          } ]
        } ]
      } ]
    }`

	// deletion would result into unassignment
	utils.DelConfig(t, jsonCfg, testbed.GetNodes()[0])
	time.Sleep(1 * time.Second)

	// confirm the default gateway in one of the containers
	output, err = node.RunCommandWithOutput("docker exec myContainer1 /sbin/ip route")
	if err != nil {
		t.Fatalf("Error - unable to get default ip route, output = '%s'", output)
	}
	if strings.Contains(output, "default via 10.240.0.254") {
		t.Fatalf("Error - able to still find the default rout after network is deleted output = '%s'", output)
	}

}
