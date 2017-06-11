// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

// +build darwin dragonfly freebsd netbsd openbsd

package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"unsafe"

	//"golang.org/x/net/bpf"
	"golang.org/x/sys/unix"
)

// ErrNoBpfDev occurs if no usable /dev/bpf* device was found.
var ErrNoBpfDev = errors.New("no usable bpf device")

type ivalue struct {
	name  [unix.IFNAMSIZ]byte
	value int16
}

func ioctl(fd int, req int, p unsafe.Pointer) (err error) {
	_, _, errno := unix.Syscall(
		unix.SYS_IOCTL,
		uintptr(fd),
		uintptr(req),
		uintptr(p),
	)
	if errno != 0 {
		err = error(errno)
	}

	return
}

func listen(ifi *net.Interface, proto uint16) (fd, l int, err error) {
	// Find a usable Berkeley Packet Filter device.
	var f *os.File
	for i := 0; ; i++ {
		dev := fmt.Sprintf("/dev/bpf%d", i)

		// bpf devices should exist in sequence from 0.. so if it
		// doesn't exist then we've exhausted all of our options.
		_, err = os.Stat(dev)
		if err != nil {
			break
		}

		// Attempt to open device.
		f, err = os.OpenFile(dev, os.O_RDWR, 0666)
		if err == nil {
			break
		}
	}
	if err != nil {
		err = ErrNoBpfDev
		return
	}
	fd = int(f.Fd())

	err = setBpfInterface(fd, ifi.Name)
	if err != nil {
		return
	}

	err = setBpfImmediate(fd)
	if err != nil {
		return
	}

	// TODO: add protocol filter

	l, err = getBpfBuflen(fd)

	return
}

func joinMcast(ifi *net.Interface, fd int, addr string) (err error) {
	// TODO: add multicast address filter
	return
}

func getBpfBuflen(fd int) (l int, err error) {
	err = ioctl(fd, unix.BIOCGBLEN, unsafe.Pointer(&l))
	return
}

func setBpfImmediate(fd int) (err error) {
	m := 1 // Enabled
	err = ioctl(fd, unix.BIOCIMMEDIATE, unsafe.Pointer(&m))
	return
}

func setBpfInterface(fd int, name string) (err error) {
	var iv ivalue
	copy(iv.name[:], []byte(name))
	err = ioctl(fd, unix.BIOCSETIF, unsafe.Pointer(&iv))
	return
}
