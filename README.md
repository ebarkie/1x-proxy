# 802.1x Port Access Entity Proxy

This proxies PAE protocol packets between interfaces. Per the 802.1d
specification this traffic is not supposed to be bridged so most
switches will not do so.

For ISP's that require 802.1x authentication it can be useful to
front-end the CPE with a device that runs this.  It has been tested
on various Ubiquiti EdgeRouter devices.

**This currently only works for Linux.  The BSD implementation is not complete.**

## Installation

Go already makes cross-compiling incredibly simple but a Makefile is
included for convenience.

* Native build
```
$ make
```

* EdgeRouter Lite/Pro build
```
$ make mips64
```

* EdgeRouter X [SFP] build
```
$ make mipsle
```

## Usage

```
Usage of ./1x-proxy:
  -names string
    	colon delimited interface list (default "eth0:eth1")
```

## License

Copyright (c) 2016-2019 Eric Barkie. All rights reserved.  
Use of this source code is governed by the MIT license
that can be found in the [LICENSE](LICENSE) file.
