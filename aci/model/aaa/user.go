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

package aaa

import (
	"github.com/contiv/netplugin/aci/model"
	"encoding/xml"
)

type User struct {
	model.Mo
	XMLName xml.Name `xml:"aaaUser"`
	Name string `xml:"name,attr,omitempty"`
	Pwd string  `xml:"pwd,attr,omitempty"`
	AccountStatus string  `xml:"accountStatus,attr,omitempty"`
	FirstName string  `xml:"firstName,attr,omitempty"`
	LastName string  `xml:"lastName,attr,omitempty"`
}