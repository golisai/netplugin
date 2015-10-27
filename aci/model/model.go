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

package model

import (
	"fmt"
	"reflect"
)

// Mo is the base struct for all the model structs
// All MO structs should embed this struct
type Mo struct {
	Dn string `xml:"dn,attr,omitempty"`
}

// MoClass returns the className of the underlying Mo
//  v can be any of pointer, slice or a struct
func MoClass(v interface{}) (string, error) {
	ty := reflect.TypeOf(v)
	for ty.Kind() == reflect.Ptr || ty.Kind() == reflect.Slice {
		ty = ty.Elem()
	}
	
	if ty.Kind() != reflect.Struct {
		return "", fmt.Errorf("%v is not an Mo struct", ty.Name())
	}
	
	fld, ok := ty.FieldByName("XMLName");
	if !ok {
		return "", fmt.Errorf("Invalid Mo struct %v - no XMLName filed defined", ty.Name())
	}

	className := fld.Tag.Get("xml")
	if className == "" {
		return "", fmt.Errorf("Invalid Mo struct %v - no xml tag for XMLName filed defined", ty.Name())
	}

	return className, nil
}
