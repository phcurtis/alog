// Copyright 2017 phcurtis grplog Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grplog_test

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/phcurtis/fn"
	"github.com/phcurtis/grplog"
)

func Example_mustnew() {
	l := grplog.MustNew("log:", grplog.FlagsDef)
	l.Info.Println("CallStk:" + fn.CStk())
	/* Representative Output of the above follows <lines wrapped>:
	log:INFO: 2017/10/22 21:43:06 example_test.go:18  FN:grplog_test.ExampleMustNew()
		CallStk:github.com/phcurtis/grplog_test.ExampleMustNew
		<--testing.runExample<--testing.runExamples<--testing.(*M).Run<--main.main
		<--runtime.main<--runtime.goexit
	*/
}

func Example_new() {
	//
	l, err := grplog.New("log:", grplog.FfnFull)
	if err != nil {
		panic(err)
	}
	l.Info.Println("CurTime:", time.Now())
}

func Example_newspecial() {
	iowr := grplog.IowrDefault()
	iowr.Error = os.Stdout

	b, err := grplog.NewSpecial("blog:", grplog.FlagsDef, 0, iowr)
	if err != nil {
		panic(err)
	}
	b.SetFlags(grplog.LflagsOff)
	b.Trace.SetOutput(ioutil.Discard)
	b.SetPkgFlags(grplog.LflagsDef)
	b.Trace.Println("this trace is discarded -- won't be seen")
	b.Alert.Println("alert newspecial b1")
	b.Error.Println("error newspecial b1")
	b.Notice.Println("notice msg1")
	b.Notice.SetIgnore(true)
	b.Notice.Println("notice msg2 this will be ignored")
	b.Notice.SetIgnore(false)
	b.Notice.Println("notice msg3")

	// Output:
	// blog:ALERT: FN:grplog_test.Example_newspecial() alert newspecial b1
	// blog:ERROR: FN:grplog_test.Example_newspecial() error newspecial b1
	// blog:NOTICE: FN:grplog_test.Example_newspecial() notice msg1
	// blog:NOTICE: FN:grplog_test.Example_newspecial() notice msg3

}

func Example_verifyLvls() {
	fn.LogSetOutput(os.Stdout)
	// This example exercises a basic test on all grplog levels
	a := grplog.MustNew("glog:", grplog.FlagsOff)
	// turn off log flags so don't have to deal with time in output string
	a.SetFlags(0)
	a.SetPkgFlags(0)
	//defer a.TraceFnEnd(a.TraceFnBeg())

	a.Trace.Println("Trace-A1")
	a.Debug.Printf("%s\n", "Debug-A1")
	a.Debug.Print("FN:"+fn.Cur(), " Debug-A2\n")
	a.Info.Println("Info-A1")
	a.Notice.Println("Notice-A1")
	a.Warning.Println("Warning-A1")
	a.Alert.Println("Alert-A1")

	a.Error.SetOutput(ioutil.Discard)
	a.Error.Println("Error-A1-Discard-will-not-see")

	a.Critical.SetOutput(ioutil.Discard)
	a.Critical.Println("Critical-A1-Discard-will-not-see")

	a.Emergency.SetOutput(ioutil.Discard)
	a.Emergency.Println("Emergency-A1-Discard-will-not-see")

	a.Error.SetOutput(os.Stdout)
	a.Error.Println("Error-A1")

	a.Critical.SetOutput(os.Stdout)
	a.Critical.Println("Critical-A1")

	a.Emergency.SetOutput(os.Stdout)
	a.Emergency.Println("Emergency-A1")

	b, err := grplog.New("blog:", 0)
	if err != nil {
		panic(err)
	}
	b.SetFlags(0)
	b.Trace.Println("TRACE-B1")

	a.SetLabel("GLOG:")
	a.Alert.Println("ALERT-LOGNAME-CHANGEA1")

	// Output:
	// glog:TRACE: Trace-A1
	// glog:DEBUG: Debug-A1
	// glog:DEBUG: FN:github.com/phcurtis/grplog_test.Example_verifyLvls Debug-A2
	// glog:INFO: Info-A1
	// glog:NOTICE: Notice-A1
	// glog:WARNING: Warning-A1
	// glog:ALERT: Alert-A1
	// glog:ERROR: Error-A1
	// glog:CRITICAL: Critical-A1
	// glog:EMERGENCY: Emergency-A1
	// blog:TRACE: TRACE-B1
	// GLOG:ALERT: ALERT-LOGNAME-CHANGEA1
}

func Example_groupprintln() {
	g, err := grplog.NewSpecial("glog:", 0, grplog.LflagsOff, grplog.IowrDefault())
	// Error, Critical and Emergency by default go to standard error so set them to stdout
	g.Error.SetOutput(os.Stdout)
	g.Critical.SetOutput(os.Stdout)
	g.Emergency.SetOutput(os.Stdout)
	if err != nil {
		log.Panic(err)
	}
	g.Println(`"this a test"`)
	// Output:
	// glog:TRACE: "this a test"
	// glog:DEBUG: "this a test"
	// glog:INFO: "this a test"
	// glog:NOTICE: "this a test"
	// glog:WARNING: "this a test"
	// glog:ALERT: "this a test"
	// glog:ERROR: "this a test"
	// glog:CRITICAL: "this a test"
	// glog:EMERGENCY: "this a test"
}

func Example_printvarious() {
	g, err := grplog.NewSpecial("glog:", grplog.FlagsOff, grplog.LflagsOff, grplog.IowrDefault())
	if err != nil {
		log.Panic(err)
	}

	// Error, Critical and Emergency by default go to standard error so set them to stdout
	g.Error.SetOutput(os.Stdout)
	g.Critical.SetOutput(os.Stdout)
	g.Emergency.SetOutput(os.Stdout)

	g.CondPrintln(true, `"this a test"`)
	g.CondPrintln(false, `"should not see"`)
	g.Trace.Println("*****traceLevel-println*****")
	g.CondPrintln(true, `"should see"`)
	g.Trace.CondPrintln(false, "shoud NOT see")
	g.Debug.CondPrintln(true, "from debug")
	g.SetIgnoreAll(true)
	g.CondPrintln(true, "should NOT see")
	g.Println("again should *NOT* see this")
	g.Critical.CondPrintln(true, "should NOT see")
	g.SetIgnoreAll(false)
	g.Critical.CondPrintln(true, "criticalMSG")
	g.Critical.Println("criticalMSG2")
	g.Info.Println("InfoMsg")
	g.Info.SetIgnore(true)
	g.Info.Println("should *NOT* see")
	g.Info.CondPrint(true, "again should *NOT* see")
	g.Notice.CondPrint(true, "should see this notice")
	g.Trace.CondPrintf(true, "one:%d two:%d res:%d\n", 1, 2, 3)
	g.Trace.CondPrintf(false, "1 Should NOT see one:%d two:%d res:%d\n", 1, 2, 3)
	g.Trace.SetIgnore(true)
	g.Trace.CondPrintf(true, "2 Should NOT see one:%d two:%d res:%d\n", 1, 2, 3)
	g.Trace.Printf("3 Should NOT see one:%d two:%d res:%d\n", 1, 2, 3)
	g.Debug.SetPkgFlags(grplog.FfnFull)
	g.Debug.Printf("<=full funcname")
	g.Debug.SetPkgFlags(grplog.FfnBase | grplog.Ffilenogps)
	g.Debug.SetFlags(log.Llongfile)
	g.Debug.Print("<=longfile less gps")
	g.Debug.SetIgnore(true)
	g.Debug.Print("should NOT see <=base funcname")
	g.Debug.Print("should not see")
	// Output:
	// glog:TRACE: "this a test"
	// glog:DEBUG: "this a test"
	// glog:INFO: "this a test"
	// glog:NOTICE: "this a test"
	// glog:WARNING: "this a test"
	// glog:ALERT: "this a test"
	// glog:ERROR: "this a test"
	// glog:CRITICAL: "this a test"
	// glog:EMERGENCY: "this a test"
	// glog:TRACE: *****traceLevel-println*****
	// glog:TRACE: "should see"
	// glog:DEBUG: "should see"
	// glog:INFO: "should see"
	// glog:NOTICE: "should see"
	// glog:WARNING: "should see"
	// glog:ALERT: "should see"
	// glog:ERROR: "should see"
	// glog:CRITICAL: "should see"
	// glog:EMERGENCY: "should see"
	// glog:DEBUG: from debug
	// glog:CRITICAL: criticalMSG
	// glog:CRITICAL: criticalMSG2
	// glog:INFO: InfoMsg
	// glog:NOTICE: should see this notice
	// glog:TRACE: one:1 two:2 res:3
	// glog:DEBUG: FN:github.com/phcurtis/grplog_test.Example_printvarious() <=full funcname
	// glog:DEBUG: github.com/phcurtis/grplog/example_test.go:185 FN:grplog_test.Example_printvarious() <=longfile less gps

}
