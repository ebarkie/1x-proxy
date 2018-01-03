// Copyright (c) 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package main

import (
	"io"
	"net"

	"golang.org/x/sys/unix"
)

// PacketConn represents a raw packet connection for an Ethernet
// interface.
type PacketConn struct {
	fd     int
	ifi    *net.Interface
	buflen int
}

// NewPacketConn creates a PacketConn using the specified Ethernet interface.
func NewPacketConn(ifi *net.Interface, proto uint16) (p *PacketConn, err error) {
	p = &PacketConn{ifi: ifi}
	p.fd, p.buflen, err = listen(ifi, proto)
	return
}

// Buflen returns the size of a PacketConn buffer.
func (p *PacketConn) Buflen() int {
	return p.buflen
}

// Close closes a PacketConn.
func (p *PacketConn) Close() error {
	return unix.Close(p.fd)
}

// JoinMcast joins an Ethernet multicast MAC address.
func (p *PacketConn) JoinMcast(addr string) (err error) {
	return joinMcast(p.ifi, p.fd, addr)
}

// Read reads from a PacketConn.
func (p *PacketConn) Read(b []byte) (n int, err error) {
	n, err = unix.Read(p.fd, b)
	if err != nil {
		n = 0
	}

	return
}

// String returns the PacketConn interface name.
func (p *PacketConn) String() string {
	return p.ifi.Name
}

// Write writes to a PacketConn.
func (p *PacketConn) Write(b []byte) (n int, err error) {
	n, err = unix.Write(p.fd, b)
	if err != nil {
		return
	}
	if n != len(b) {
		err = io.ErrShortWrite
	}

	return
}
