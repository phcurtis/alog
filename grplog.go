// Copyright 2017 phcurtis grplog Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package grplog creates a grouped set of convenience hooks into stdlib log.logger
// for various levels of logging: Trace, Debug, Info, Notice, Warning, Alert, Error,
// Critical, Emergency. Having them grouped easily allows multiple sets within a
// program as well as providing a way to easily grep the output coming from that
// group.  i.e.  glog:TRACE: glog:DEBUG versus blog:TRACE: blog:ERROR  ....
package grplog

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// Version of this package
var Version = 0.220

//type flagst int

// func name type setting during log output
const (
	FfnBase    = 1 << iota            // no funcname output
	FfnFull                           // base funcname output
	Ffilenogps                        // remove go path src prefix applicable to FfnFull
	FlagsDef   = FfnBase | Ffilenogps //
	FlagsOff   = 0
)

// stdlib log flags convenience constants
const (
	LflagsDTS  = log.Ldate | log.Ltime | log.Lshortfile
	LflagsDTL  = log.Ldate | log.Ltime | log.Llongfile
	LflagsDTSM = log.Ldate | log.Ltime | log.Lshortfile | log.Lmicroseconds
	LflagsDTLM = log.Ldate | log.Ltime | log.Llongfile | log.Lmicroseconds
	LflagsDef  = LflagsDTS
	LflagsDefL = LflagsDTL
	LflagsOff  = 0
)

var gopathsrc string

func init() {
	gopathsrc = os.Getenv("GOPATH")
	if gopathsrc != "" {
		gopathsrc += "/src/"
	}
}

// Log related constants
const (
	LogAlignFileDef = 24 // log alignment 'file' field minimum width
	LogAlignFileMax = 50 // log alignment 'file' field minimum width max

	LogAlignFuncDef = 0  // log alignment 'func' field minimum width
	LogAlignFuncMax = 50 // log alignment 'func' field minimum width max
)

type alignStruct struct {
	filea int
	funca int
}

// Gtrace - global trace log which is available at init time of this package.
var Gtrace = LvlStruct{
	flags:     FlagsDef,
	log:       log.New(os.Stdout, "GTRACE: ", LflagsDef),
	mu:        new(sync.Mutex),
	logOutput: os.Stdout,
	name:      "Gtrace",
	align:     alignStruct{filea: LogAlignFileDef, funca: LogAlignFuncDef},
}

// LvlStruct - contains a given log level stuff
type LvlStruct struct {
	ignore     bool        // way to ignore Print,Printf,Println, CondPrint, CondPrintln
	flags      int         // func name type
	log        *log.Logger // stdlib logger
	par        *GlvlStruct // parent this lvl belongs too if nil its Gtrace
	mu         *sync.Mutex //
	logOutput  io.Writer   // maintain copy since log.logger does not support Get Output
	outCtr     uint64      // counter of times func 'out' called
	outCharCtr uint64      // counter of chars sent through func 'out' and onto log.logger
	name       string      // go entryPoint name
	align      alignStruct //
}

// GlvlStruct - group log level struct
type GlvlStruct struct {
	ignoreall    bool // way to ignore Print,Printf,Println, CondPrint, CondPrintln
	Name         string
	Trace        *LvlStruct
	Debug        *LvlStruct
	Info         *LvlStruct
	Notice       *LvlStruct
	Warning      *LvlStruct
	Alert        *LvlStruct
	Error        *LvlStruct
	Critical     *LvlStruct
	Emergency    *LvlStruct
	mu           sync.Mutex // mutex for group
	firstIowr    IowrStruct
	logAlignFile int
	logAlignFunc int
}

// IowrStruct - grplog iowriters struct
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

// IowrDefault returns the current default IowrStruct. Subsequent
// invocations could be different if io.Writers such as os.Stdout
// and os.Stderr have changed between invocations of this function.
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

// Base labels for each log level
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

type lvlListStruct struct {
	level **LvlStruct
	name  string
	Blab  string
	iowr  *io.Writer
}

func (g *GlvlStruct) lvlList() []lvlListStruct {
	return []lvlListStruct{
		{&g.Trace, "Trace", TraceBlab, &g.firstIowr.Trace},
		{&g.Debug, "Debug", DebugBlab, &g.firstIowr.Debug},
		{&g.Info, "Info", InfoBlab, &g.firstIowr.Info},
		{&g.Notice, "Notice", NoticeBlab, &g.firstIowr.Notice},
		{&g.Warning, "Warning", WarningBlab, &g.firstIowr.Warning},
		{&g.Alert, "Alert", AlertBlab, &g.firstIowr.Alert},
		{&g.Error, "Error", ErrorBlab, &g.firstIowr.Error},
		{&g.Critical, "Critical", CriticalBlab, &g.firstIowr.Critical},
		{&g.Emergency, "Emergency", EmergencyBlab, &g.firstIowr.Emergency},
	}
}

// NewSpecial ... returns *GlvlStruct and error based on following arguments.
//	- glabel is grplog group label applied to all levels of logging
//		and concatenated with each specific level logger base labels [Blab]
//  - flags - func name type
//	- logFlagsGroup is stdlib log - log flags value to apply to all levels of logging
// 	- iowr is to contain corresponding iowriters for all logging levels
//		[one could use ioutil.Discard to inactivate a level of logging]
func NewSpecial(glabel string, flags int, logFlagsGroup int, iowr IowrStruct) (*GlvlStruct, error) {
	return newll(glabel, flags, logFlagsGroup, &iowr, false)
}

// New ... returns *GlvlStruct and error based on following arguments.
// see NewSpecial on input parameters.
func New(glabel string, flags int) (*GlvlStruct, error) {
	return newll(glabel, flags, LflagsDef, nil, false)
}

// MustNew returns *GlvlStruct and panics if any error occurs
// see NewSpecial on input parameters
func MustNew(glabel string, flags int) *GlvlStruct {
	g, _ := newll(glabel, flags, LflagsDef, nil, true)
	return g
}

var grplogCount uint32
var muGrplogCount sync.Mutex

// newll - worker func that creates a new blogStruct
func newll(glabel string, flags int, logFlags int, iowr *IowrStruct, panicErr bool) (*GlvlStruct, error) {
	g := &GlvlStruct{firstIowr: IowrDefault()}
	if iowr != nil {
		g.firstIowr = *iowr
	}
	// test g.firstIowr.Error = nil

	muGrplogCount.Lock()
	grplogCount++
	g.Name = fmt.Sprintf("%s<%d>", glabel, grplogCount) //ensure an unique name
	muGrplogCount.Unlock()

	for _, v := range g.lvlList() {
		if *v.iowr == nil {
			errnew := errors.New(v.name + " io.Writer is nil, if you want to discard use ioutil.Discard")
			if panicErr {
				log.Panic(errnew)
			}
			return nil, errnew
		}
		*v.level = &LvlStruct{
			log:       log.New(*v.iowr, glabel+v.Blab, logFlags),
			logOutput: *v.iowr,
			flags:     flags,
			mu:        &g.mu,
			par:       g,
			name:      v.name,
			align:     alignStruct{filea: LogAlignFileDef, funca: LogAlignFuncDef},
		}
	}
	return g, nil
}
