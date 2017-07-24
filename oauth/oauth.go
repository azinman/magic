// Copyright Aaron Zinman 2017, 2018
// Copyright Duck Research LLC 2017, 2018
// All rights reserved.
//
// This file is part of Magichaus.
//
// Magichaus is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Magichaus is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Magichaus.  If not, see <http://www.gnu.org/licenses/>.

package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type OAuth struct {
	clientId        string
	clientSecret    string
	deviceCode      string
	UserCode        string
	VerificationUrl string
	expiresAt       time.Time
	pollInterval    int
}

func NewGoogleOAuth(clientId string, clientSecret string) (*OAuth, error) {
	resp, err := http.PostForm("https://accounts.google.com/o/oauth2/device/code",
		url.Values{"client_id": {clientId}, "scope": {"email profile"}})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	res := struct {
		device_code      string
		user_code        string
		verification_url string
		expires_in       int
		interval         int
	}{}
	err = decoder.Decode(&res)
	if err != nil {
		return nil, err
	}
	return &OAuth{
		clientId:        clientId,
		clientSecret:    clientSecret,
		deviceCode:      res.device_code,
		UserCode:        res.user_code,
		VerificationUrl: res.verification_url,
		expiresAt:       time.Now().Add(time.Duration(res.expires_in) * time.Second),
		pollInterval:    res.interval,
	}, nil
}

func (o *OAuth) NextStepString() string {
	return fmt.Sprintf("Go to %v and enter the code %v", o.VerificationUrl, o.UserCode)
}

func (o *OAuth) Poll() {
	// needt os tart polling here in a go routine + chan
	// https://developers.google.com/identity/sign-in/devices#get_user_profile_information_from_the_id_token
}
