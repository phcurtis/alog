/* package alog creates a grouped set of convenience hooks into stdlib log.logger
 for various levels of logging: Trace, Debug, Info, Notice, Warning, Alert, Error,
 Critical, Emergency. Having them grouped easily allows multiple sets within a 
 program as well as providing a way to easily grep the output coming from that 
 group.  i.e.  alog:TRACE: alog:DEBUG versus blog:TRACE: blog:ERROR  ....
*/
package alog

import (
	"errors"
	"io"
	"log"
	"os"
)

var Version = "0.04"

// alog level base labels [Blab]
const (
	TraceBlab     = "TRACE: "
	DebugBlab     = "DEBUG: "
	InfoBlab      = "INFO: "
	NoticeBlab    = "NOTICE: "
	WarningBlab   = "WARNING: "
	AlertBlab     = "ALERT: "
	ErrorBlab     = "ERROR: "
	CriticalBlab  = "CRITICAL: "
	EmergencyBlab = "EMERGENCY: "
)

// stdlib log flags
const (
	LFlagsDef = log.Ldate |
		log.Ltime |
		log.Lmicroseconds |
		log.Lshortfile

	LFlagsDefL = log.Ldate |
		log.Ltime |
		log.Lmicroseconds |
		log.Llongfile

	LFlagsCmn = log.Ldate |
		log.Ltime |
		log.Lshortfile

	LFlagsCmnL = log.Ldate |
		log.Ltime |
		log.Llongfile
)

// alog iowriters struct
type IowrStruct struct {
	Trace     io.Writer
	Debug     io.Writer
	Info      io.Writer
	Notice    io.Writer
	Warning   io.Writer
	Alert     io.Writer
	Error     io.Writer
	Critical  io.Writer
	Emergency io.Writer
}

// IowrDefault returns the current default IowrStruct which subsequent invocations could be different if
// io.Writers such as os.Stdout and os.Stderr have changed between invocations of this function.
func IowrDefault() IowrStruct {
	return IowrStruct{

		// default to os.Stdout
		Trace:   os.Stdout,
		Debug:   os.Stdout,
		Info:    os.Stdout,
		Notice:  os.Stdout,
		Warning: os.Stdout,
		Alert:   os.Stdout,

		// default to os.Stderr
		Error:     os.Stderr,
		Critical:  os.Stderr,
		Emergency: os.Stderr,
	}
}

// alog struct containing corresponding distinct loggers for each level
type AlogStruct struct {
	Trace     *log.Logger
	Debug     *log.Logger
	Info      *log.Logger
	Notice    *log.Logger
	Warning   *log.Logger
	Alert     *log.Logger
	Error     *log.Logger
	Critical  *log.Logger
	Emergency *log.Logger
	firstIowr IowrStruct //contains original values for iowriters when AlogStruct was created
	// individual alog level io.Writers can be changed after initialization via .SetOutput and also if referencing things like os.Stdout could be altered via some os.Stdout type of redirection ...
}

// firstIowr returns the IowrStruct contents on a specific AlogStruct creation
func (a *AlogStruct) FirstIowr() IowrStruct {
	return a.firstIowr
}

// SetGroupLogFlags - sets stdlib log flags for all alog levels to fval.
// fval - stdlib log flag value.
// Calling this func with fval=0 can be handy during testing (especially to avoid dealing with [changing] time output).
// If you have other goroutines that may also call this or explicitly such as a.Trace.SetFlag() you are
// responsible in coordinating that possibly with a mutex.
func (a *AlogStruct) SetGroupLogFlags(fval int) {
	for _, v := range a.alogLBList() {
		v.Level.SetFlags(fval)
	}
}

// SetGroupLabel - applies alog group string to all alog levels.
// glabel - alog group label
// If you have other goroutines that may also call this or explicitly such as a.Trace.SetPrefix() you are
// responsible in coordinating that possibly with a mutex
func (a *AlogStruct) SetGroupLabel(glabel string) {
	for _, v := range a.alogLBList() {
		v.Level.SetPrefix(glabel + v.Blab)
	}
}

/* NewSpecial returns *AlogStruct and error based on following arguments.
 	- iowr is to contain corresponding iowriters for all logging levels
		[one could use ioutil.Discard to inactivate a level of logging]
 	- glabel is alog group label applied to all levels of logging
		and concatenated with each specific level logger base labels [Blab]
	- logFlagsGroup is stdlib log - log flags value to apply to all levels of logging
*/
func NewSpecial(iowr IowrStruct, glabel string, logFlagsGroup int) (*AlogStruct, error) {
	return newll(&iowr, glabel, logFlagsGroup)
}

/* New returns *AlogStruct and error
 	- glabel is alog group label applied to all levels of logging
		and concatenated with each specific level logger base labels
*/
func New(glabel string) (*AlogStruct, error) {
	return newll(nil, glabel, LFlagsDef)
}

/*
******** NON public and supporting items ********
 */

// alog (L)evel (B)aselabel List struct
type alogLBListStruct struct {
	Level *log.Logger // alog level
	Blab  string      // alog level base label
}

func (a *AlogStruct) alogLBList() []alogLBListStruct {
	return []alogLBListStruct{
		{a.Trace, TraceBlab},
		{a.Debug, DebugBlab},
		{a.Info, InfoBlab},
		{a.Notice, NoticeBlab},
		{a.Warning, WarningBlab},
		{a.Alert, AlertBlab},
		{a.Error, ErrorBlab},
		{a.Critical, CriticalBlab},
		{a.Emergency, EmergencyBlab},
	}
}

// newll creates entry points for convenience logger hooks calling back into stdlib log.Logging
func newll(iowr *IowrStruct, glabel string, logFlags int) (*AlogStruct, error) {
	alog := AlogStruct{firstIowr: IowrDefault()}
	if iowr != nil {
		alog.firstIowr = *iowr
	}
	if alog.firstIowr.Trace == nil ||
		alog.firstIowr.Debug == nil ||
		alog.firstIowr.Info == nil ||
		alog.firstIowr.Notice == nil ||
		alog.firstIowr.Warning == nil ||
		alog.firstIowr.Alert == nil ||
		alog.firstIowr.Error == nil ||
		alog.firstIowr.Critical == nil ||
		alog.firstIowr.Emergency == nil {
		return nil, errors.New("one or more elements of iowrStruct are nil for io.Writer, you might consider using ioutil.Discard for such elements")
	}

	alog.Trace = log.New(alog.firstIowr.Trace, glabel+TraceBlab, logFlags)
	alog.Debug = log.New(alog.firstIowr.Debug, glabel+DebugBlab, logFlags)
	alog.Info = log.New(alog.firstIowr.Info, glabel+InfoBlab, logFlags)
	alog.Notice = log.New(alog.firstIowr.Notice, glabel+NoticeBlab, logFlags)
	alog.Warning = log.New(alog.firstIowr.Warning, glabel+WarningBlab, logFlags)
	alog.Alert = log.New(alog.firstIowr.Alert, glabel+AlertBlab, logFlags)
	alog.Error = log.New(alog.firstIowr.Error, glabel+ErrorBlab, logFlags)
	alog.Critical = log.New(alog.firstIowr.Critical, glabel+CriticalBlab, logFlags)
	alog.Emergency = log.New(alog.firstIowr.Emergency, glabel+EmergencyBlab, logFlags)

	return &alog, nil
}
