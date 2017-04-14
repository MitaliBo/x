// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The go1.4.3 command runs the go command from go1.4.3.
//
// To install, run:
//
//     $ go get golang.org/x/build/version/go1.4.3
//     $ go1.4.3 download
//
// And then use the go1.4.3 command as if it were your normal go
// command.
//
// See the release notes at https://beta.golang.org/doc/go1.4
//
// File bugs at http://golang.org/issues/new
package main

import "golang.org/x/build/version"

func main() {
	version.Run("go1.4.3")
}