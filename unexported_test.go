// Copyright 2017 phcurtis grplog Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grplog

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

func Test_osExit(t *testing.T) {
	want := "exitcode=1"
	var osExitStr string
	osExitSave := osExit
	osExit = func(code int) {
		osExitStr = fmt.Sprint("exitcode=", code)
	}
	if !testing.Verbose() {
		log.SetOutput(ioutil.Discard)
	}
	defer func() {
		if !testing.Verbose() {
			log.SetOutput(os.Stderr)
		}
		osExit = osExitSave
	}()

	g := MustNew("glog:", 0)
	if !testing.Verbose() {
		g.Trace.SetOutput(ioutil.Discard)
	}
	g.Trace.Fatal("invoked Trace.Fatal cause os.Exit(1)")
	if osExitStr != want {
		t.Errorf("Trace.Fatal got:%q want:%q", osExitStr, want)
	}
	osExitStr = ""
	g.Trace.Fatalf("invoked Trace.Fatalf cause os.Exit(1)")
	if osExitStr != want {
		t.Errorf("Trace.Fatalf got:%q want:%q", osExitStr, want)
	}
	osExitStr = ""
	g.Trace.Fatalln("invoked Trace.Fatalln cause os.Exit(1)")
	if osExitStr != want {
		t.Errorf("Trace.Fatalln got:%q want:%q", osExitStr, want)
	}
}

func Test_newllpanic(t *testing.T) {
	if !testing.Verbose() {
		log.SetOutput(ioutil.Discard) // toss log.Panic output
	}
	defer log.SetOutput(os.Stderr) // restore log output

	defer func() {
		var p interface{}
		p = recover()
		if testing.Verbose() {
			log.Printf("panicErr:%v\n", p)
		}
		if p == nil {
			t.Errorf("should have paniced ... due to passing in nil io.Writer")
		}
		if testing.Verbose() {
			log.Println("Recovered from panic")
		}
	}()
	iowr := IowrDefault()
	iowr.Debug = nil
	_, _ = newll("glog:", 0, LflagsDef, &iowr, true)
}

func Test_align(t *testing.T) {
	tests := []struct {
		name  string
		width int
		str   string
		want  string
	}{
		{"t1", 20, "1234567890", "1234567890" + strings.Repeat(" ", 10)},
		{"t2", 10, "1234567890", "1234567890"},
		{"t3", 4, "4321-3", "4321-3"},
		{"t4", 5, "1", "1    "},
	}
	for _, test := range tests {
		input := test.str
		got := align(input, test.width)
		if got != test.want {
			t.Errorf("align(%q,%d): got:%q want:%q\n", input, test.width, got, test.want)
		}
	}
}

func Test_groupsetflags(t *testing.T) {
	g := MustNew("glog:", 0)
	input := []int{0, 0xffff, 0xdead}
	want := []int{0, 0xffff, 0xdead}

	for i, v := range input {
		g.SetFlags(v)
		for _, v1 := range g.lvlList() {
			got := (*v1.level).log.Flags()
			if got != want[i] {
				t.Errorf("%s.SetFlags(0x%x) got:0x%x want:0x%x\n", v1.name, v, got, want[i])
			}
		}
	}
}

func Test_groupsetignore(t *testing.T) {
	g := MustNew("glog:", 0)
	input := []bool{true, false, true, false}
	want := []bool{true, false, true, false}

	for i, v := range input {
		g.SetIgnore(v)
		for _, v1 := range g.lvlList() {
			got := (*v1.level).Ignore()
			if got != want[i] {
				t.Errorf("%s.SetIgnore(%t) got:%t want:%t\n", v1.name, v, got, want[i])
			}
		}
	}
}

func Test_groupsetlabel(t *testing.T) {
	g := MustNew("glog:", 0)
	input := []string{"a", " ", "grplog:", "AlOg337z:"}
	want := []string{"a", " ", "grplog:", "AlOg337z:"}

	for i, v := range input {
		g.SetLabel(v)
		for _, v1 := range g.lvlList() {
			got := (*v1.level).Prefix()
			if got != want[i]+v1.Blab {
				t.Errorf("%s.SetLabel(%s) got:%s want:%s\n", v1.name, v, got, want[i])
			}
		}
	}
}

func Test_groupsetpkgflags(t *testing.T) {
	g := MustNew("glog:", 0)
	input := []int{0, 0x1234, 0xffffffff}
	want := []int{0, 0x1234, 0xffffffff}

	for i, v := range input {
		g.SetPkgFlags(v)
		for _, v1 := range g.lvlList() {
			got := (*v1.level).PkgFlags()
			if got != want[i] {
				t.Errorf("%s.SetPkgFlags(0x%x) got:0x%x want:0x%x\n", v1.name, v, got, want[i])
			}
		}
	}
}

func Test_groupsetoutput(t *testing.T) {
	g := MustNew("glog:", 0)
	input := []io.Writer{os.Stderr, os.Stdout}
	want := []io.Writer{os.Stderr, os.Stdout}

	for i, v := range input {
		g.SetOutput(v)
		for _, v1 := range g.lvlList() {
			got := (*v1.level).GetOutput()
			if got != want[i] {
				t.Errorf("%s.SetOutput(%s) got:%s want:%s\n", v1.name, v, got, want[i])
			}
		}
	}
}
