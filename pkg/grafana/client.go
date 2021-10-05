//--------------------------------------------------------------------------
// Copyright 2018 infinimesh
// www.infinimesh.io
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
//--------------------------------------------------------------------------

package grafana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Client struct {
	c       *http.Client
	baseURL string

	user     string
	password string
}

func NewClient(baseURL string, user, password string) *Client {
	return &Client{
		c:        &http.Client{},
		baseURL:  baseURL,
		user:     user,
		password: password,
	}
}

func (c *Client) CreateUser(name string) error {
	js, err := json.Marshal(map[string]string{
		"name":     name,
		"login":    name,
		"password": "dummy",
	})
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(js)
	req, err := http.NewRequest("POST", c.baseURL+"/api/admin/users", buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	req.SetBasicAuth(c.user, c.password)

	resp, err := c.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Wrong status code: %v", resp.StatusCode)
	}
	return nil
}

func (c *Client) CreateOrg(name string) error {
	js, err := json.Marshal(map[string]string{
		"name": name,
	})
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(js)
	req, err := http.NewRequest("POST", c.baseURL+"/api/orgs", buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	req.SetBasicAuth(c.user, c.password)
	resp, err := c.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Wrong status code: %v", resp.StatusCode)
	}
	return nil

}

func (c *Client) AddUserToOrg(orgID int, name string, role string) error {
	js, err := json.Marshal(map[string]string{
		"loginOrEmail": name,
		"role":         role,
	})
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(js)
	req, err := http.NewRequest("POST", c.baseURL+"/api/orgs/"+strconv.Itoa(orgID)+"/users", buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	req.SetBasicAuth(c.user, c.password)
	resp, err := c.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Wrong status code: %v", resp.StatusCode)
	}
	return nil
}

func (c *Client) GetOrgID(orgName string) (orgID int, err error) {

	type Org struct {
		Id int `json:"id"`
	}

	req, err := http.NewRequest("GET", c.baseURL+"/api/orgs/name/"+orgName, http.NoBody)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	req.SetBasicAuth(c.user, c.password)
	resp, err := c.c.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	buf := &bytes.Buffer{}
	var org Org

	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal(buf.Bytes(), &org)
	if err != nil {
		return 0, err
	}

	return org.Id, nil

}

func (c *Client) GetUserID(username string) (userID int, err error) {
	type User struct {
		Id int `json:"id"`
	}

	req, err := http.NewRequest("GET", c.baseURL+"/api/users/lookup?loginOrEmail="+username, http.NoBody)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	req.SetBasicAuth(c.user, c.password)
	resp, err := c.c.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	buf := &bytes.Buffer{}
	var user User

	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal(buf.Bytes(), &user)
	if err != nil {
		return 0, err
	}

	return user.Id, nil

}

func (c *Client) MakeUserAdmin(userID int) error {
	js, err := json.Marshal(map[string]interface{}{
		"isGrafanaAdmin": true,
	})
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(js)
	req, err := http.NewRequest("PUT", c.baseURL+"/api/admin/users/"+strconv.Itoa(userID)+"/permissions", buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	req.SetBasicAuth(c.user, c.password)
	resp, err := c.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Wrong status code: %v", resp.StatusCode)
	}
	return nil
}

func (c *Client) SwitchUserOrg(userID int, orgID int) error {
	req, err := http.NewRequest("POST", c.baseURL+"/api/users/"+strconv.Itoa(userID)+"/using/"+strconv.Itoa(orgID), http.NoBody)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	req.SetBasicAuth(c.user, c.password)
	resp, err := c.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Wrong status code: %v", resp.StatusCode)
	}
	return nil

}
