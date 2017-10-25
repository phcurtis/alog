// Copyright 2017 phcurtis grplog Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grplog_test

import (
	"testing"

	"github.com/phcurtis/grplog"
)

func TestSetIgnoreAll(t *testing.T) {
	g := grplog.MustNew("glog:", 0)
	input := []bool{true, false, true, false}
	want := []bool{true, false, true, false}
	for i, _ := range input {
		g.SetIgnoreAll(input[i])
		got := g.GetIgnoreAll()
		if got != want[i] {
			t.Errorf("[%s].SetIgnore(%t) got:%t want:%t\n", g.Name, input[i], got, want[i])
		}
	}
}
