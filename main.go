// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package main

//go:generate ./version.sh

import (
	"flag"
	"log"
	"strings"
)

func main() {
	// Read interface names from flag arguments.
	names := flag.String("names", "eth0:eth1", "colon delimited interface list")
	flag.Parse()

	log.Printf("802.1x Port Access Entity Proxy (version %s)", version)

	// Add listeners.
	for _, name := range strings.Split(*names, ":") {
		err := AddProxyMember(name)
		if err != nil {
			log.Fatalf("Error setting up interface %s: %s", name, err.Error())
		}
	}

	// Block forever.
	select {}
}
