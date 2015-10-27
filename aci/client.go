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
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/contiv/netplugin/aci/model"
	"github.com/contiv/netplugin/aci/model/aaa"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
)

var _ = log.Print

type MoNotFoundError struct {
	Dn string
}

func (e MoNotFoundError) Error() string {
	return fmt.Sprintf("Dn %v not found", e.Dn)
}


type RestError struct {
	XMLName xml.Name `xml:"error"`
	Code string `xml:"code,attr"`
	Text string `xml:"text,attr"`
}

func (e RestError) Error() string {
	return fmt.Sprintf("Rest error %v: %v", e.Code, e.Text)
}

// Client stores the ACI Rest access info 
type Client struct {
	http.Client
	URL string
	User string
	Passwd string
}

// NewClient creates new ACI client
func NewClient(url string, user string, passwd string) (*Client, error) {
	jar, err :=  cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	c := &Client{http.Client{Jar: jar}, url, user, passwd}
	err = c.Login()
	return c, err
}

// Login posts the ACI login policy
func (c *Client) Login() error {
	user := &aaa.User{Name: c.User, Pwd: c.Passwd}
	return c.Post(c.URL + "/api/aaaLogin.xml", user)
}

func (c *Client) Post(url string, v interface{}) error {
	data, err := xml.Marshal(v)
	if err != nil {
		return err
	}
	
	resp, err := c.Client.Post(url, "appication/xml", bytes.NewReader(data))
	if err != nil {
		return err
	}
	resp.Body.Close()


	return nil
}

// Query sends the rest query and unmarshals the response
func (c *Client) query(url string, className string, v interface{}) (count uint32, err error) {
	// Send the query string
	resp, err := c.Client.Get(url)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return
	}

	// Decode the response, skipping the imdata tag
	decoder := xml.NewDecoder(bytes.NewReader(body))
	var t xml.Token
	for t, err = decoder.Token(); err == nil; t, err = decoder.Token() {
		if elem, ok := t.(xml.StartElement); ok {
			switch elem.Name.Local {
			case "error":
				var e RestError
				if err = decoder.DecodeElement(&e, &elem); err == nil {
					err = e
				}	
			case className:
				if err = decoder.DecodeElement(v, &elem); err == nil {
					count++
				}
			}

			if err != nil {
				break
			}
		}
	}

	if err == io.EOF {
		err = nil
	}

	return
}

// ClassQuery sends MO class query and unmarshals the response
func (c *Client) ClassQuery(v interface{}) error {
	className, err := model.MoClass(v)
	if err != nil {
		return err
	}
	url := c.URL + "/api/node/class/" + className + ".xml" 
	_, err = c.query(url, className, v)
	return err
}

// DnQuery sends the DN query and unmarshals the response
func (c *Client) DnQuery(dn string, v interface{}) (error) {
	className, err := model.MoClass(v)
	if err != nil {
		return err
	}
	url := c.URL + "/api/node/mo/" + dn + ".xml"
	count, err := c.query(url, className, v)
	if err == nil && count == 0 {
		err = MoNotFoundError{dn}
	}
	return err
}
