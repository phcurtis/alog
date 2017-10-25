// Copyright 2017 phcurtis grplog Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grplog_test

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/phcurtis/grplog"
)

// support function for main bench marking
func brtn(b *testing.B, db dbs) {
	if testing.Verbose() {
		log.Printf("testname:%s b.N is %d\n", b.Name(), b.N)
	}
	a, err := grplog.New(db.logname, 0)
	if err != nil {
		b.Fatal(err)
	}
	a.SetIgnoreAll(db.ignoreall)

	var iowr io.Writer
	var fnamestr string
	if db.tempsfn != "" {
		tmpfile, err := ioutil.TempFile("", "Benchmark-"+db.tempsfn)
		if err != nil {
			b.Fatal(err)
		}
		defer func() {
			if testing.Verbose() {
				log.Printf("\nremoving tempfile:%s\n", tmpfile.Name())
			}
			err := os.Remove(tmpfile.Name())
			if err != nil {
				log.Printf("\nerror removing %v err:%v\n", tmpfile.Name(), err)
			}
		}()
		//defer tmpfile.Close() maybe TODO to allow files to linger so can study output after bench mark
		fnamestr = " Routing log output to tempfile:" + tmpfile.Name()
		iowr = tmpfile
	} else {
		iowr = ioutil.Discard
		fnamestr = " Routing log output to ioutil.Discard"
	}

	if testing.Verbose() {
		log.Println(fnamestr)
	}

	b.ResetTimer()

	if db.alllevels {
		a.SetFlags(db.lflags)
		a.SetOutput(iowr)
	} else {
		a.Trace.SetFlags(db.lflags)
		a.Trace.SetOutput(iowr)
	}

	for i := 0; i < b.N; i++ {
		for i := 0; i < db.reps; i++ {
			if db.alllevels {
				a.Println("produced via GroupPrintln")
			} else {
				a.Trace.Println("blab:", grplog.TraceBlab)
			}
		}
	}
}

type dbs struct {
	ignoreall bool   //
	alllevels bool   // test all log levels
	reps      int    // repetitions for each level
	logname   string // group name
	tempsfn   string // temp suffix filename if empty string then route to ioutil.Discard
	lflags    int    // stdlib log flags
}

func BenchmarkAlog(b *testing.B) {

	log.SetOutput(os.Stdout)
	b.Run("lfdef-s1x.", func(b *testing.B) {
		brtn(b, dbs{false, false, 1, "glog1:", "test1-", grplog.LflagsDef})
	})

	b.Run("lfdef-s1x.-ignoreall", func(b *testing.B) {
		brtn(b, dbs{true, false, 1, "glog1:", "test1ign-", grplog.LflagsDef})
	})

	b.Run("lfdefL-s1x.", func(b *testing.B) {
		brtn(b, dbs{false, false, 1, "glog2:", "test2-", grplog.LflagsDefL})
	})
	b.Run("lfdtsm-s1x", func(b *testing.B) {
		brtn(b, dbs{false, false, 1, "glog3:", "test3-", grplog.LflagsDTSM})
	})
	b.Run("lfdtlm-s1x", func(b *testing.B) {
		brtn(b, dbs{false, false, 1, "glog4:", "test4-", grplog.LflagsDTLM})
	})

	b.Run("lfdtlm-s1x-ignoreall", func(b *testing.B) {
		brtn(b, dbs{true, false, 1, "glog4:", "test4-", grplog.LflagsDTLM})
	})

	b.Run("lfoff-s1x.", func(b *testing.B) {
		brtn(b, dbs{false, false, 1, "glog5:", "test5-", grplog.LflagsOff})
	})

	b.Run("lfshortf-s1x", func(b *testing.B) {
		brtn(b, dbs{false, false, 1, "glog6:", "", log.Lshortfile})
	})
	b.Run("lflongf-s1x", func(b *testing.B) {
		brtn(b, dbs{false, false, 1, "glog7:", "", log.Llongfile})
	})
	b.Run("lfoff-s1x", func(b *testing.B) {
		brtn(b, dbs{false, false, 1, "glog8:", "", 0})
	})
	b.Run("lfoff-m1x", func(b *testing.B) {
		brtn(b, dbs{false, true, 1, "glog8:", "", 0})
	})
	b.Run("lfoff-m10x", func(b *testing.B) {
		brtn(b, dbs{false, true, 10, "glog8:", "", 0})
	})

	b.Run("lfoff-m10x-ignoreAll", func(b *testing.B) {
		brtn(b, dbs{true, true, 10, "glog8:", "", 0})
	})
}
