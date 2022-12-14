package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func readRespBody(resp *http.Response) []byte {
	if resp.StatusCode != 200 {
		fmt.Println("status:", resp.Status)
		return nil
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return data
}

// ZtClient is a wrapper of general operation for zerotier
type ZtClient struct {
	client  *http.Client
	timeout time.Duration
	baseURL string
	token   string
	uid     string

	display func(o interface{})
}

// Init setup the backend http client and init the members
func (c *ZtClient) Init() {
	if c.timeout == 0 {
		c.timeout = 20 * time.Second
	}

	c.client = &http.Client{Timeout: c.timeout}
	c.baseURL = strings.TrimSuffix(c.baseURL, "/")
}

// NewRequest assemble http request with token and possible payload
func (c *ZtClient) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, err
	}

	var req *http.Request
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		req, err = http.NewRequest(method, u.String(), nil)
		if err != nil {
			return nil, err
		}
	default:
		buf := new(bytes.Buffer)
		if body != nil {
			err = json.NewEncoder(buf).Encode(body)
			if err != nil {
				return nil, err
			}
		}

		req, err = http.NewRequest(method, u.String(), buf)
		if err != nil {
			return nil, err
		}
	}

	// request with token
	req.Header.Add("Authorization", "token "+c.token)

	return req, nil
}

// DoRequest send http request and parse the response
func (c *ZtClient) DoRequest(method, path string, body, data interface{}) (*http.Response, error) {
	// request
	req, err := c.NewRequest(method, path, body)
	if err != nil {
		fmt.Println("** Failed to construct request")
		return nil, nil
	}

	resp, err := c.client.Do(req)
	if err != nil {
		fmt.Println("** error:", err)
		return nil, nil
	}

	// handle
	if method == http.MethodPost {
		defer resp.Body.Close()
	}

	payload := readRespBody(resp)
	if data != nil {
		err := json.Unmarshal(payload, &data)
		if err != nil {
			fmt.Println("** error:", err)
		}
	}

	return resp, err
}

// GetUIDHack get uid bind to the token by create a immediate network
func (c *ZtClient) GetUIDHack() string {
	// get uid by create a tmp network
	body := make(map[string]interface{})
	data := make(map[string]interface{})
	c.DoRequest(http.MethodPost, "/network", &body, &data)

	// once get resp, remove the network
	if nid, ok := data["id"]; ok {
		if v, ok := nid.(string); ok {
			c.DoRequest(http.MethodDelete, "/network/"+v, nil, nil)
		}
	}

	// parse and return uid
	if uid, ok := data["ownerId"]; ok {
		if v, ok := uid.(string); ok {
			return v
		}
	}

	return ""
}

// DumpUserRecord dumps user record info identified by the token
func (c *ZtClient) DumpUserRecord() {
	// get uid if not specified
	if c.uid == "" {
		c.uid = c.GetUIDHack()
	}

	if c.uid == "" {
		fmt.Println("** Failed to get UID")
		return
	}

	// get user record
	data := make(map[string]interface{})
	c.DoRequest(http.MethodGet, "/user/"+c.uid, nil, &data)
	c.display(data)
}

// ListNetwork list all managed networks
func (c *ZtClient) ListNetwork() {
	data := make([]ZtNetInfo, 0)
	c.DoRequest(http.MethodGet, "/network", nil, &data)

	c.display(data)
}

// CreateNetwork create new network and apply the the settings
func (c *ZtClient) CreateNetwork(conf ZtNetPost) {
	body := make(map[string]interface{})
	data := make(map[string]interface{})
	c.DoRequest(http.MethodPost, "/network", &body, &data)

	if nid, ok := data["id"]; ok {
		if v, ok := nid.(string); ok {
			c.SetNetwork(v, conf)
		}
	}
}

// SetNetwork apply new setting to existing network identified by nid
func (c *ZtClient) SetNetwork(nid string, conf ZtNetPost) {
	data := ZtNetInfo{}
	c.DoRequest(http.MethodPost, "/network/"+nid, &conf, &data)

	c.display(data)
}

// DelNetwork delete a network and print the response status
func (c *ZtClient) DelNetwork(nid string) {
	resp, err := c.DoRequest(http.MethodDelete, "/network/"+nid, nil, nil)
	if err == nil && resp.StatusCode == 200 {
		fmt.Println("status:", resp.Status)
	}
}

// ListNetworkMember list members of a/all network
func (c *ZtClient) ListNetworkMember(nid string) {
	if nid == "" {
		data := make([]ZtNetInfo, 0)
		c.DoRequest(http.MethodGet, "/network", nil, &data)

		for i := range data {
			c.ListNetworkMember(data[i].ID)
		}
	} else {
		path := "/network/" + nid + "/member"
		data := make([]ZtNetMemberInfo, 0)
		c.DoRequest(http.MethodGet, path, nil, &data)

		fmt.Printf("-- net: %s\n", nid)
		c.display(data)
		fmt.Println("")
	}
}

// SetNetworkMember apply settings like name and ip to the network member
func (c *ZtClient) SetNetworkMember(nid string, mid string, conf ZtNetMemberPost) {
	path := "/network/" + nid + "/member/" + mid
	data := ZtNetMemberInfo{}

	c.DoRequest(http.MethodPost, path, &conf, &data)
	c.display(data)
}

// DelNetworkMember delete member of a network and print the response status
func (c *ZtClient) DelNetworkMember(nid string, mid string) {
	path := "/network/" + nid + "/member/" + mid
	resp, err := c.DoRequest(http.MethodDelete, path, nil, nil)
	if err == nil && resp.StatusCode == 200 {
		fmt.Println("status:", resp.Status)
	}
}
