// Copyright 2017 phcurtis grplog Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grplog_test

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/phcurtis/grplog"
)

func panicRecover(p interface{}, panictext string) error {
	if testing.Verbose() {
		log.Printf("panicErr:%v\n", p)
	}
	if p == nil {
		return errors.New("should have paniced ... " + panictext)
	}
	return nil
}

func Test_levelpanic(t *testing.T) {
	name := "Trace.Panic"
	defer func() {
		if err := panicRecover(recover(), "due to calling "+name); err != nil {
			t.Error(err)
		} else if testing.Verbose() {
			log.Println("Recovered from " + name)
		}
	}()
	g := grplog.MustNew("glog:", 0)
	if !testing.Verbose() {
		g.Trace.SetOutput(ioutil.Discard)
	}
	defer hideStderr(t)()
	g.Trace.Panic("Trace.Panic called")
	t.Error(name + " level did NOT throw panic")
}

func Test_levelpanicf(t *testing.T) {
	name := "Trace.Panicf"
	defer func() {
		if err := panicRecover(recover(), "due to calling "+name); err != nil {
			t.Error(err)
		} else if testing.Verbose() {
			log.Println("Recovered from " + name)
		}
	}()
	g := grplog.MustNew("glog:", 0)
	if !testing.Verbose() {
		g.Trace.SetOutput(ioutil.Discard)
	}
	defer hideStderr(t)()
	g.Trace.Panicf("Trace.Panicf called")
	t.Error(name + " level did NOT throw panic")
}

func Test_levelpanicln(t *testing.T) {
	name := "Trace.Panicln"
	defer func() {
		if err := panicRecover(recover(), "due to calling "+name); err != nil {
			t.Error(err)
		} else if testing.Verbose() {
			log.Println("Recovered from " + name)
		}
	}()
	g := grplog.MustNew("glog:", 0)
	if !testing.Verbose() {
		g.Trace.SetOutput(ioutil.Discard)
	}

	defer hideStderr(t)()
	g.Trace.Panicln("Trace.Panicln called")
	t.Error(name + " level did NOT throw panic")
}

func hideStderr(t *testing.T) func() {
	tmpfile, err := ioutil.TempFile("", "stderrCap-")
	if err != nil {
		panic(err)
	}
	if !testing.Verbose() {
		log.SetOutput(ioutil.Discard)
	}
	stderrSave := os.Stderr
	os.Stderr = tmpfile
	return func() {
		if !testing.Verbose() {
			log.SetOutput(os.Stderr)
		}
		os.Stderr = stderrSave
		_ = os.Remove(tmpfile.Name())
	}
}
