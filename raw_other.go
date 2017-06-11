// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

// +build !darwin,!dragonfly,!freebsd,!linux,!netbsd,!openbsd

package main

import (
	"errors"
	"net"
)

// ErrNotImplemented occurs if this is run on an Operating System
// where the lower level calls have not yet been implemented.
var ErrNotImplemented = errors.New("not implemented")

func listen(ifi *net.Interface, proto uint16) (fd, l int, err error) {
	err = ErrNotImplemented
	return
}

func joinMcast(ifi *net.Interface, fd int, addr string) (err error) {
	err = ErrNotImplemented
	return
}
