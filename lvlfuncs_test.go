// Copyright 2017 phcurtis grplog Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grplog_test

import (
	"io"
	"os"
	"testing"

	"github.com/phcurtis/grplog"
)

func Test_lvlfuncs(t *testing.T) {

	l := grplog.Gtrace

	tests := []struct {
		name string
		fset interface{}
		arg1 interface{}
		fget interface{}
		arg2 interface{}
	}{

		{"AlignFile",
			l.SetAlignFile, []int{1, 10, -2, grplog.LogAlignFileMax + 1, grplog.LogAlignFileMax},
			l.AlignFile, []int{1, 10, 0, grplog.LogAlignFileMax, grplog.LogAlignFileMax}},
		{"AlignFunc",
			l.SetAlignFunc, []int{1, 10, -2, grplog.LogAlignFuncMax + 1, grplog.LogAlignFuncMax},
			l.AlignFunc, []int{1, 10, 0, grplog.LogAlignFuncMax, grplog.LogAlignFuncMax}},
		{"Flags",
			l.SetFlags, []int{0xdeadbeef, 0},
			l.Flags, []int{0xdeadbeef, 0}},
		{"Ignore",
			l.SetIgnore, []bool{true, false, true},
			l.Ignore, []bool{true, false, true}},
		{"Output",
			l.SetOutput, []io.Writer{os.Stderr, os.Stdout, os.Stderr},
			l.GetOutput, []io.Writer{os.Stderr, os.Stdout, os.Stderr}},
		{"Prefix",
			l.SetPrefix, []string{"xiferP", "pReFiX", "12345:"},
			l.Prefix, []string{"xiferP", "pReFiX", "12345:"}},
		{"PkgFlags",
			l.SetPkgFlags, []int{0xdeadbeef, 0, 0xffffffff},
			l.PkgFlags, []int{0xdeadbeef, 0, 0xffffffff}},
	}
	for _, v := range tests {
		switch fset := v.fset.(type) {
		case func(int):
			islice := v.arg1.([]int)
			wslice := v.arg2.([]int)
			fget := v.fget.(func() int)
			for i := 0; i < len(islice); i++ {
				input := islice[i]
				fset(input)
				got := fget()
				want := wslice[i]
				if got != want {
					t.Errorf("%s(%v) got:%v want:%v\n", v.name, input, got, want)
				}
			}
		case func(bool):
			islice := v.arg1.([]bool)
			wslice := v.arg2.([]bool)
			fget := v.fget.(func() bool)
			for i := 0; i < len(islice); i++ {
				input := islice[i]
				fset(input)
				got := fget()
				want := wslice[i]
				if got != want {
					t.Errorf("%s(%v) got:%v want:%v\n", v.name, input, got, want)
				}
			}
		case func(io.Writer):
			islice := v.arg1.([]io.Writer)
			wslice := v.arg2.([]io.Writer)
			fget := v.fget.(func() io.Writer)
			for i := 0; i < len(islice); i++ {
				input := islice[i]
				fset(input)
				got := fget()
				want := wslice[i]
				if got != want {
					t.Errorf("%s(0x%x) got:0x%x want:0x%x\n", v.name, input, got, want)
				}
			}
		case func(string):
			islice := v.arg1.([]string)
			wslice := v.arg2.([]string)
			fget := v.fget.(func() string)
			for i := 0; i < len(islice); i++ {
				input := islice[i]
				fset(input)
				got := fget()
				want := wslice[i]
				if got != want {
					t.Errorf("%s(%v) got:%v want:%v\n", v.name, input, got, want)
				}
			}
		default:
			t.Errorf("test: %q not supported", v.name)
		}
	}
}
