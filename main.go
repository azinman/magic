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

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/azinman/magic/alexa"
	"github.com/azinman/magic/housedb"
	"github.com/azinman/magic/tv"
	"github.com/azinman/magic/wifi"
	// _ "github.com/azinman/magic/pipeline"
	// _ "github.com/azinman/magic/sinks"
	// _ "github.com/azinman/magic/sources"
)

var (
	flagTV = flag.Bool("tv", true, "Connect to TV")
)

type Magic struct {
	tv     *tv.TV
	db     housedb.Db
	httpDb *housedb.HTTPDB
	alexa  *alexa.Alexa
}

func Alexa() *alexa.Alexa {
	return magic.alexa
}

func TV() *tv.TV {
	return magic.tv
}

func HTTPDB() *housedb.HTTPDB {
	return magic.httpDb
}

func DB() housedb.Db {
	return magic.db
}

func newMagic() *Magic {
	log.Println("Magic starting...")
	db, err := housedb.NewMemory()
	if err != nil {
		log.Fatal("Couldn't load HouseDB: ", err)
	}
	// hdb, err := housedb.NewHTTPDB(db)
	// if err != nil {
	// 	log.Fatal("Couldn't load HouseDB: ", err)
	// }

	// alexa, err := alexa.New()
	// if err != nil {
	// 	log.Fatal("Couldn't load Alexa: ", err)
	// }
	// var t *tv.TV
	// if *flagTV {
	// 	var err error
	// 	t, err = tv.NewTV()
	// 	if err != nil {
	// 		log.Fatal("Couldn't connect to tv:", err)
	// 	}
	// 	device.Register(t)
	// }
	// if err := dns.AWSUpdateToPublicIP(dns.AWSZoneID, "755.magicha.us", 60); err != nil {
	// 	log.Fatal("Couldn't update IP", err)
	// }

	for e := range wifi.MonitorWifi() {
		fmt.Println(e)
	}

	return &Magic{alexa: nil, db: db, tv: nil, httpDb: nil}
}

var magic *Magic

func init() {
	flag.Parse()
}

func main() {
	magic = newMagic()
	// log.Print("Magic inited; starting HTTP server")
	// HTTPDB().Run()
}
