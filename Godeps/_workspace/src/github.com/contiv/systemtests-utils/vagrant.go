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

package utils

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
)

// Vagrant implements a vagrant based testbed
type Vagrant struct {
	expectedNodes int
	nodes         map[string]TestbedNode
}

// Setup brings up a vagrant testbed
func (v *Vagrant) Setup(start bool, env string, numNodes int) error {
	v.nodes = map[string]TestbedNode{}

	vCmd := &VagrantCommand{ContivNodes: numNodes, ContivEnv: env}

	if start {
		output, err := vCmd.RunWithOutput("up")
		if err != nil {
			log.Errorf("Vagrant up failed. Error: %s Output: \n%s\n",
				err, output)
			return err
		}

		defer func() {
			if err != nil {
				v.Teardown()
			}
		}()
	}

	v.expectedNodes = numNodes

	output, err := vCmd.RunWithOutput("status")
	if err != nil {
		log.Errorf("Vagrant status failed. Error: %s Output: \n%s\n",
			err, output)
		return err
	}

	// now some hardwork of finding the names of the running nodes from status output
	re, err := regexp.Compile("[a-zA-Z0-9_\\- ]*running \\(virtualbox\\)")
	if err != nil {
		return err
	}
	nodeNamesBytes := re.FindAll(output, -1)
	if nodeNamesBytes == nil {
		err = fmt.Errorf("No running nodes found in vagrant status output: \n%s\n",
			output)
		return err
	}
	nodeNames := []string{}
	for _, nodeNameByte := range nodeNamesBytes {
		nodeName := strings.Fields(string(nodeNameByte))[0]
		nodeNames = append(nodeNames, nodeName)
	}
	if len(nodeNames) != numNodes {
		err = fmt.Errorf("Number of running node(s) (%d) is not equal to number of expected node(s) (%d) in vagrant status output: \n%s\n",
			len(nodeNames), numNodes, output)
		return err
	}

	// some more work to figure the ssh port and private key details
	// XXX: vagrant ssh-config --host <> seems to be broken as-in it doesn't
	// correctly filter the output based on passed host-name. So filtering
	// the output ourselves below.
	if output, err = vCmd.RunWithOutput("ssh-config"); err != nil {
		return fmt.Errorf("Error running vagrant ssh-config. Error: %s. Output: \n%s\n", err, output)
	}

	if re, err = regexp.Compile("Host [a-zA-Z0-9_-]+|Port [0-9]+|IdentityFile .*"); err != nil {
		return err
	}

	nodeInfosBytes := re.FindAll(output, -1)
	if nodeInfosBytes == nil {
		return fmt.Errorf("Failed to find node info in vagrant ssh-config output: \n%s\n", output)
	}

	// got the names, now fill up the vagrant-nodes structure
	for _, nodeName := range nodeNames {
		nodeInfoPos := -1
		// nodeInfos is a slice of node info orgranised as nodename, it's port and identity-file in that order per node
		for j := range nodeInfosBytes {
			if string(nodeInfosBytes[j]) == fmt.Sprintf("Host %s", nodeName) {
				nodeInfoPos = j
				break
			}
		}
		if nodeInfoPos == -1 {
			return fmt.Errorf("Failed to find %q info in vagrant ssh-config output: \n%s\n", nodeName, output)
		}
		port := ""
		if n, err := fmt.Sscanf(string(nodeInfosBytes[nodeInfoPos+1]), "Port %s", &port); n == 0 || err != nil {
			return fmt.Errorf("Failed to find %q port info in vagrant ssh-config output: \n%s\n. Error: %s",
				nodeName, nodeInfosBytes[nodeInfoPos+1], err)
		}
		privKeyFile := ""
		if n, err := fmt.Sscanf(string(nodeInfosBytes[nodeInfoPos+2]), "IdentityFile %s", &privKeyFile); n == 0 || err != nil {
			return fmt.Errorf("Failed to find %q identity file info in vagrant ssh-config output: \n%s\n. Error: %s",
				nodeName, nodeInfosBytes[nodeInfoPos+2], err)
		}
		log.Infof("Adding node: %q(%s:%s)", nodeName, port, privKeyFile)
		var node *VagrantNode
		if node, err = NewVagrantNode(nodeName, port, privKeyFile); err != nil {
			return err
		}
		v.nodes[node.GetName()] = TestbedNode(node)
	}

	return nil
}

// Teardown cleans up a vagrant testbed
func (v *Vagrant) Teardown() {
	for _, node := range v.nodes {
		vnode := node.(*VagrantNode)
		vnode.Cleanup()
	}
	vCmd := &VagrantCommand{ContivNodes: v.expectedNodes}
	output, err := vCmd.RunWithOutput("destroy", "-f")
	if err != nil {
		log.Errorf("Vagrant destroy failed. Error: %s Output: \n%s\n",
			err, output)
	}

	v.nodes = map[string]TestbedNode{}
	v.expectedNodes = 0
}

// GetNode obtains a node by name.
func (v *Vagrant) GetNode(name string) TestbedNode {
	return v.nodes[name]
}

// GetNodes returns the nodes in a vagrant setup
func (v *Vagrant) GetNodes() []TestbedNode {
	var ret []TestbedNode
	for _, value := range v.nodes {
		ret = append(ret, value)
	}

	return ret
}

// IterateNodes walks each host and executes the function supplied. On error,
// it waits for all hosts to complete before returning the error.
func (v *Vagrant) IterateNodes(fn func(TestbedNode) error) error {
	wg := sync.WaitGroup{}
	nodes := v.GetNodes()
	errChan := make(chan error, len(nodes))

	for _, node := range nodes {
		wg.Add(1)

		go func(node TestbedNode) {
			if err := fn(node); err != nil {
				errChan <- fmt.Errorf(`Error: "%v" on host: %q"`, err, node.GetName())
			}
			wg.Done()
		}(node)
	}

	wg.Wait()

	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}

// SSHAllNodes will ssh into each host and run the specified command.
func (v *Vagrant) SSHAllNodes(cmd string) error {
	return v.IterateNodes(func(node TestbedNode) error {
		return node.RunCommand(cmd)
	})
}
