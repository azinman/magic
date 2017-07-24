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

package alexa

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

const redirectURL = "https://pitangui.amazon.com/api/skill/link/XXXXXXXX"

var (
	flagSSLPEMKeyPath = flag.String("alexaSSLPEMKeyPath", "", "Path to SSL Private key in PEM")
	// flagSSLKeyP12Password = flag.String("alexaSSLKeyP12Password", "", "Password to SSL Private p12-encoded key")
	flagSSLPEMCertPath = flag.String("alexaSSLPEMCertPath", "", "Path to SSL Certificate in PEM")
	flagAlexaHTTPPort  = flag.Int("alexaHTTPPort", 4001, "Port that alexa sends to")
	// flagSSLCACertPath     = flag.String("keySSLCACertPath", "", "Path to SSL CA Certificate in PEM")
)

// Alexa is the main structure that houses the connection to Alexa
type Alexa struct {
	router *gin.Engine
}

// New creates a new Alexa
func New() (*Alexa, error) {
	if *flagSSLPEMKeyPath == "" {
		return nil, errors.New("-alexaSSLPEMKeyPath not defined")
	}
	if *flagSSLPEMCertPath == "" {
		return nil, errors.New("-alexaSSLPEMCertPath not defined")
	}
	r := newRouter()
	addr := fmt.Sprintf(":%d", *flagAlexaHTTPPort)
	certFile := *flagSSLPEMCertPath
	keyFile := *flagSSLPEMKeyPath
	if err := r.RunTLS(addr, certFile, keyFile); err != nil {
		return nil, err
	}
	return &Alexa{router: r}, nil
}

// HTTPS -----------------------------------------------------------------------

func newRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/alexaOAuthLogin", alexaOAuthLogin)
	r.GET("/alexaOAuthToken", alexaOAuthToken)
	return r
}

// CERTIFICATES ----------------------------------------------------------------

// func loadP12(path string, password string) ([]byte, error) {
// 	b, err := ioutil.ReadFile(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	pemBlocks, err := pkcs12.ToPEM(b, password)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var pemData []byte
// 	for _, block := range pemBlocks {
// 		pemData = append(pemData, pem.EncodeToMemory(block)...)
// 	}
// 	return pemData, nil
// }
//
// func loadCerts() ([]tls.Certificate, error) {
// 	privatePEM, err := loadP12(*flagSSLKeyP12Path, *flagSSLKeyP12Password)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to load P12: %v", err)
// 	}
// 	log.Print("Loaded private pem: ", privatePEM)
// 	// Cert key
// 	publicPEM, err := ioutil.ReadFile(*flagSSLCertPath)
// 	if err != nil {
// 		return nil, err
// 	}
// 	cert, err := tls.X509KeyPair(publicPEM, privatePEM)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return []tls.Certificate{cert}, nil
// }

// HANDLERS --------------------------------------------------------------------

func alexaOAuthLogin(c *gin.Context) {
	// log.Printf("Got oauth request=%v params=%v keys=%v", c.Request, c.Params, c.Keys)
	clientID := c.Query("client_id")
	if clientID != "magic-haus-alexa-skill" {
		log.Print("Don't understand HTTP request: ", c.Request)
		c.AbortWithStatus(404)
		return
	}
	state := c.Query("state")
	responseType := c.Query("response_type")
	scope := c.Query("scope")
	if state == "" || responseType != "code" || scope != "ALEXA_SCOPE" {
		log.Print("Don't understand HTTP request: ", c.Request)
		c.AbortWithStatus(500)
		return
	}
	log.Printf("Got oauth: state=%v client_id=%v response_type=%v scope=%v",
		state, clientID, responseType, scope)
	code := "DEADBEEF"
	url, _ := url.Parse(redirectURL)
	q := url.Query()
	q.Set("state", state)
	q.Set("code", code)
	url.RawQuery = q.Encode()
	log.Print("Redirecting to ", url.String())
	c.Redirect(http.StatusTemporaryRedirect, url.String())
}

func alexaOAuthToken(c *gin.Context) {
	log.Print("Got", c.Request.URL.Query())
	c.AbortWithStatus(500)
	return
}
