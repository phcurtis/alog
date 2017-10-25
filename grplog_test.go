package alog

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestBasicNoOutput(t *testing.T) {
	a, err := New("alog:")
	if err != nil {
		t.Fatal(err)
	}
	iowr := a.FirstIowr()
	iowr.Trace = nil
	_, err = NewSpecial(iowr, "blog:", 0)
	if err == nil {
		t.Fatal("should have failed since iowr.Trace was nil")
	}
}

// support function for main bench marking
func brtn(db dbs, b *testing.B) {
	a, err := New(db.logname)
	if err != nil {
		b.Fatal(err)
	}

	var iowr io.Writer
	var fnamestr string
	if db.tempbfn != "" {
		tmpfile, err := ioutil.TempFile("", db.tempbfn)
		if err != nil {
			b.Fatal(err)
		}
		defer func() {
			if testing.Verbose() {
				fmt.Println("removing tempfile:" + tmpfile.Name())
				os.Remove(tmpfile.Name())
			}
		}()
		//defer tmpfile.Close() maybe TODO to allow files to linger so can study output after bench mark
		fnamestr = "Routing log output to tempfile:" + tmpfile.Name()
		iowr = tmpfile
	} else {
		iowr = ioutil.Discard
		fnamestr = "Routing log output to ioutil.Discard"
	}
	if testing.Verbose() {
		fmt.Println(fnamestr)
	}

	b.ResetTimer()

	if db.alllevels {
		for _, v := range a.alogLBList() {
			v.Level.SetFlags(db.lflags)
			v.Level.SetOutput(iowr)
		}
	} else {
		a.Trace.SetFlags(db.lflags)
		a.Trace.SetOutput(iowr)
	}

	for i := 0; i < b.N; i++ {
		for i := 0; i < db.reps; i++ {
			if db.alllevels {
				for _, v := range a.alogLBList() {
					v.Level.Println("blab:", v.Blab)
				}
			} else {
				a.Trace.Println("blab:", TraceBlab)
			}
		}
	}
}

type dbs struct {
	alllevels bool   // test all log levels
	reps      int    // repetitions for each level
	logname   string // alog name
	tempbfn   string // temp base filename if empty string then route to ioutil.Discard
	lflags    int    // stdlib log flags
}

/*  bench mark as of 2017 Oct shows if you don't want a given logging level recorded
set it to ioutil.Discard AND set log flags to zero which likely will reduce overhead a bunch
without removing the call from the code. And later you could arm it as needed with hooks to your code
calling corresponding .SetOutput and .SetFlags, etc.
*/

func BenchmarkAlog(b *testing.B) {
	b.Run("lfdef-1x", func(b *testing.B) { brtn(dbs{false, 1, "alog1:", "test1-", LFlagsDef}, b) })
	b.Run("lfdefL1x", func(b *testing.B) { brtn(dbs{false, 1, "alog2:", "test2-", LFlagsDefL}, b) })
	b.Run("lfcmn-1x", func(b *testing.B) { brtn(dbs{false, 1, "alog3:", "test3-", LFlagsCmn}, b) })
	b.Run("lfcmnL1x", func(b *testing.B) { brtn(dbs{false, 1, "alog4:", "test4-", LFlagsCmnL}, b) })
	b.Run("lfoff-1x", func(b *testing.B) { brtn(dbs{false, 1, "alog5:", "test5-", 0}, b) })
	b.Run("lfcmnD1x", func(b *testing.B) { brtn(dbs{false, 1, "alog6:", "", LFlagsCmn}, b) })
	b.Run("lfoffD1x", func(b *testing.B) { brtn(dbs{false, 1, "alog7:", "", 0}, b) })

	b.Run("lfdef-10x", func(b *testing.B) { brtn(dbs{false, 10, "alog8:", "test8-", LFlagsDef}, b) })
	b.Run("lfdefL10x", func(b *testing.B) { brtn(dbs{false, 10, "alog9:", "test9-", LFlagsDefL}, b) })
	b.Run("lfcmn-10x", func(b *testing.B) { brtn(dbs{false, 10, "alog10:", "test10-", LFlagsCmn}, b) })
	b.Run("lfcmnL10x", func(b *testing.B) { brtn(dbs{false, 10, "alog11:", "test11-", LFlagsCmnL}, b) })
	b.Run("lfoff-10x", func(b *testing.B) { brtn(dbs{false, 10, "alog12:", "test12-", 0}, b) })
	b.Run("lfcmnD10x", func(b *testing.B) { brtn(dbs{false, 10, "alog13:", "", LFlagsCmn}, b) })
	b.Run("lfoffD10x", func(b *testing.B) { brtn(dbs{false, 10, "alog14:", "", 0}, b) })
}

func Example_verifyingAllAlogLevels() {
	a, err := New("alog:")
	if err != nil {
		panic(err)
	}
	// turn off log flags so don't have to deal with time in output string
	a.SetGroupLogFlags(0)

	a.Trace.Println("Trace-A1")
	a.Debug.Println("Debug-A1")
	a.Info.Println("Info-A1")
	a.Notice.Println("Notice-A1")
	a.Warning.Println("Warning-A1")
	a.Alert.Println("Alert-A1")

	a.Error.SetOutput(ioutil.Discard)
	a.Error.Println("Error-A1-Discard")

	a.Critical.SetOutput(ioutil.Discard)
	a.Critical.Println("Critical-A1-Discard")

	a.Emergency.SetOutput(ioutil.Discard)
	a.Emergency.Println("Emergency-A1-Discard")

	a.Error.SetOutput(os.Stdout)
	a.Error.Println("Error-A1")

	a.Critical.SetOutput(os.Stdout)
	a.Critical.Println("Critical-A1")

	a.Emergency.SetOutput(os.Stdout)
	a.Emergency.Println("Emergency-A1")

	b, err := New("blog:")
	if err != nil {
		panic(err)
	}
	b.SetGroupLogFlags(0)
	b.Trace.Println("TRACE-B1")

	a.SetGroupLabel("ALOG:")
	a.Alert.Println("ALERT-LOGNAME-CHANGEA1")

	// Output:
	// alog:TRACE: Trace-A1
	// alog:DEBUG: Debug-A1
	// alog:INFO: Info-A1
	// alog:NOTICE: Notice-A1
	// alog:WARNING: Warning-A1
	// alog:ALERT: Alert-A1
	// alog:ERROR: Error-A1
	// alog:CRITICAL: Critical-A1
	// alog:EMERGENCY: Emergency-A1
	// blog:TRACE: TRACE-B1
	// ALOG:ALERT: ALERT-LOGNAME-CHANGEA1
}

func ExampleNewSpecial() {
	iowr := IowrDefault()
	iowr.Error = os.Stdout
	b, err := NewSpecial(iowr, "blog:", 0)
	if err != nil {
		panic(err)
	}
	b.Trace.SetOutput(ioutil.Discard)
	b.Trace.Println("this trace is discarded -- won't be seen")
	b.Alert.Println("alert newspecial b1")
	b.Error.Println("error newspecial b1")

	// Output:
	// blog:ALERT: alert newspecial b1
	// blog:ERROR: error newspecial b1
}
