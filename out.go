// Copyright 2017 phcurtis grplog Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grplog

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/phcurtis/fn"
)

var osExit = os.Exit

func (l *LvlStruct) outExit(s string) {
	_ = l.out(s)
	osExit(1)
}

func (l *LvlStruct) outPanic(s string) {
	_ = l.out(s)
	log.Panic(s)
}

func (l *LvlStruct) out(s string) error {
	// may have to re-examine having this lock in place for entire func
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.outll(1, s)
}

func align(str string, width int) string {
	len1 := len(str)
	if len1 < width {
		str = str + strings.Repeat(" ", width-len1)
	}
	return str
}

// out - a worker func that does final prep before calling stdlib log.Output.
func (l *LvlStruct) outll(lvladj int, s string) error {
	l.outCtr++
	//fmt.Printf("%s:outCtr:%d\n", l.name, l.outCtr)

	var fns string
	switch {
	case l.flags&FfnBase > 0:
		fns = "FN:" + fn.LvlBase(2+lvladj) + "() "
	case l.flags&FfnFull > 0:
		fns = "FN:" + fn.Lvl(2+lvladj) + "() "
	default:
	}

	l.outCharCtr += uint64(len(fns) + len(s))
	//fmt.Printf("%s:outCharCtr:%d\n", l.name, l.outCharCtr)

	lvl := 2 + lvladj
	var filenlr string

	// get original [current] log flags
	orgflags := l.log.Flags()
	sl := log.Lshortfile | log.Llongfile
	lfn := orgflags & sl

	// if log flags are including filename
	if lfn > 0 {
		_, file, line, _ := runtime.Caller(lvl)
		linenum := fmt.Sprintf(":%d", line)
		if orgflags&log.Lshortfile > 0 {
			file = filepath.Base(file)
		} else {
			// log.Llongfile
			if l.flags&Ffilenogps > 0 {
				if strings.HasPrefix(file, gopathsrc) {
					file = file[len(gopathsrc):]
				}
			}
		}

		filenlr = file + linenum + " "
		filenlr = align(filenlr, l.align.filea)

		// set log flags not to include filename
		l.log.SetFlags(orgflags &^ sl)
	}
	//logt.Printf("%s%s", filenlr, msg)

	// as of go 1.9 ... stdlib log does not check err,
	// we do here and you can decide to panic by so configuring a given level.
	err := l.log.Output(3, filenlr+fns+s)
	if lfn > 0 {
		// restore log flags
		l.log.SetFlags(orgflags)
	}
	return err
}
