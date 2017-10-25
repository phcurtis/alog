// Copyright 2017 phcurtis grplog Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grplog_test

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/phcurtis/grplog"
)

func Test_newspecialiowriternil(t *testing.T) {
	iowr := grplog.IowrDefault()
	iowr.Info = nil
	if !testing.Verbose() {
		log.SetOutput(ioutil.Discard)
		defer log.SetOutput(os.Stderr)
	}
	_, err := grplog.NewSpecial("glog", 0, grplog.LflagsDef, iowr)
	if err == nil {
		t.Errorf("should have failed since iowr.Info was nil")
	}
}
