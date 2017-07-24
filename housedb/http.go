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

package housedb

import (
	"errors"
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
)

var (
	flagHTTPPort = flag.Int("httpPort", 4000, "HTTP Port")
)

type (
	HTTPDB struct {
		router *gin.Engine
		db     Db
	}
)

var (
	flagSSLPEMKeyPath  = flag.String("SSLPEMKeyPath", "", "Path to SSL Private key in PEM")
	flagSSLPEMCertPath = flag.String("SSLPEMCertPath", "", "Path to SSL Certificate in PEM")
)

// NewHTTPDB creates a new HTTPDB
func NewHTTPDB(db Db) (*HTTPDB, error) {
	if *flagSSLPEMKeyPath == "" {
		return nil, errors.New("-SSLPEMKeyPath not defined")
	}
	if *flagSSLPEMCertPath == "" {
		return nil, errors.New("-SSLPEMCertPath not defined")
	}
	hdb := &HTTPDB{db: db}
	r := newRouter(hdb)
	hdb.router = r
	return hdb, nil
}

func (hdb *HTTPDB) Run() error {
	addr := fmt.Sprintf(":%d", *flagHTTPPort)
	certFile := *flagSSLPEMCertPath
	keyFile := *flagSSLPEMKeyPath
	return hdb.router.RunTLS(addr, certFile, keyFile)
}

// HTTPS -----------------------------------------------------------------------

func newRouter(hdb *HTTPDB) *gin.Engine {
	r := gin.Default()
	r.POST("/put", hdb.onPut)
	// 	r.GET("/alexaOAuthToken", alexaOAuthToken)
	// Get(key string) (string, error)
	// Put(key string, value string) error
	// Watch(key string, callback func(value string)) error
	// Close() error
	return r
}

// HANDLERS --------------------------------------------------------------------

func (hdb *HTTPDB) onPut(c *gin.Context) {
	key := c.PostForm("key")
	value := c.PostForm("value")
	if key == "" || value == "" || len(key) > 1024 || len(value) > 1024 {
		c.AbortWithStatus(500)
	}
	if err := hdb.db.Put(key, value); err != nil {
		c.AbortWithError(500, err)
	}
}

// func alexaOAuthToken(c *gin.Context) {
// 	log.Print("Got", c.Request.URL.Query())
// 	c.AbortWithStatus(500)
// 	return
// }
