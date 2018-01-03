// Copyright (c) 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package main

import (
	"encoding/hex"
	"errors"
	"log"
	"net"
	"sync"
)

// Errors.
var (
	ErrDupName = errors.New("Same interface name specified multiple times")
)

// PAEProxy is used to track the interfaces that are part of the proxy bridge
// and the incoming PAE packets from any of them.
type PAEProxy struct {
	members      []*PacketConn  // Active proxy members
	sync.RWMutex                // For locking members slice
	in           chan PAEPacket // Incoming PAE packets
}

// PAEPacket represents a raw PAE packet and the interface it originated
// from.
type PAEPacket struct {
	ifi *net.Interface // Interface the packet originated from
	b   []byte         // Packet bytes
}

var proxy PAEProxy

func init() {
	proxy.in = make(chan PAEPacket)

	// Run writer in the background.  This waits for new PAEPackets on the in
	// channel and for each one that arrives it writes it to each proxy member.
	go func() {
		for p := range proxy.in {
			log.Printf("%d bytes received on %s\n%s", len(p.b), p.ifi.Name, hex.Dump(p.b))

			proxy.RLock()
			for _, m := range proxy.members {
				// Skip originating interface.
				if m.ifi.Index == p.ifi.Index {
					continue
				}

				_, err := m.Write(p.b)
				if err != nil {
					log.Printf("Error sending on %s: %s", m, err.Error())
				} else {
					log.Printf("%d bytes sent on %s", len(p.b), m)
				}
			}
			proxy.RUnlock()
		}
	}()
}

// AddProxyMember adds the specified interface to the proxy bridge.
func AddProxyMember(name string) error {
	proto := uint16(0x888e) // ETH_P_PAE
	// Since syscall's are being used we're responsible for converting
	// the protocol from host byte order to network byte order.
	if !isBigEndian {
		proto = (proto<<8)&0xff00 | proto>>8
	}

	proxy.Lock()
	defer proxy.Unlock()

	// Look up interface from the name and make sure we didn't already
	// start a listener for it.
	ifi, err := net.InterfaceByName(name)
	if err != nil {
		return err
	}
	for _, m := range proxy.members {
		if ifi.Index == m.ifi.Index {
			return ErrDupName
		}
	}

	// Setup raw socket and join the multicast group that PAE packets get
	// sent to.
	m, err := NewPacketConn(ifi, proto)
	if err != nil {
		return err
	}
	err = m.JoinMcast("01:80:c2:00:00:03") // Nearest non-TPMR Bridge Group address
	if err != nil {
		m.Close()
		return err
	}

	// Run a background reader for this interface.
	proxy.members = append(proxy.members, m)
	go func() {
		log.Printf("Listening on %s", m)
		b := make([]byte, m.Buflen())
		for {
			n, err := m.Read(b)
			if err != nil {
				log.Printf("Error reading on %s: %s", m, err.Error())
				continue
			}
			proxy.in <- PAEPacket{ifi: m.ifi, b: b[:n]}
		}
	}()

	return nil
}
