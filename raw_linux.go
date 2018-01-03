// Copyright (c) 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

// +build linux

package main

import (
	"net"
	"unsafe"

	"golang.org/x/sys/unix"
)

// packetMreq is a packet membership request packet.  It's used
// to enter promiscuous mode or to add or drop membership of a
// multicast group.
type packetMreq struct {
	Ifindex int32
	Type    uint16
	Alen    uint16
	Address [8]byte
}

// sizeofPacketMreq is the size of a PacketMreq struct.
const sizeofPacketMreq = 0x10

func setSockOpt(fd, level, name int, v unsafe.Pointer, l uint32) (err error) {
	_, _, errno := unix.Syscall6(
		unix.SYS_SETSOCKOPT,
		uintptr(fd),
		uintptr(level),
		uintptr(name),
		uintptr(v),
		uintptr(l),
		0,
	)
	if errno != 0 {
		err = error(errno)
	}

	return
}

func listen(ifi *net.Interface, proto uint16) (fd, l int, err error) {
	fd, err = unix.Socket(unix.AF_PACKET, unix.SOCK_RAW, int(proto))
	if err != nil {
		return
	}

	err = unix.Bind(
		fd,
		&unix.SockaddrLinklayer{
			Protocol: proto,
			Ifindex:  ifi.Index,
		})

	l = 1522 // Max Ethernet IEEE 802.3 frame size

	return
}

func joinMcast(ifi *net.Interface, fd int, addr string) (err error) {
	mreq := packetMreq{
		Ifindex: int32(ifi.Index),
		Type:    unix.PACKET_MR_MULTICAST,
		Alen:    6, // ETH_ALEN
		Address: [8]byte{},
	}

	var a net.HardwareAddr
	a, err = net.ParseMAC(addr)
	copy(mreq.Address[0:6], a[:])
	if err != nil {
		return
	}

	err = setSockOptPacketMreq(
		fd,
		unix.SOL_PACKET,
		unix.PACKET_ADD_MEMBERSHIP,
		&mreq)

	return
}

// setSockOptPacketMreq sets packet membership for a socket.
func setSockOptPacketMreq(fd, level, opt int, mreq *packetMreq) (err error) {
	return setSockOpt(fd, level, opt, unsafe.Pointer(mreq), sizeofPacketMreq)
}
