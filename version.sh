#!/bin/sh

VERSION=`git describe --always --long --dirty`
 
cat <<EOF > version.go
package main

var version = "$VERSION"
EOF
